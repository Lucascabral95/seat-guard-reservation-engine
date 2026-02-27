package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "booking-service/docs"
	"booking-service/internal/config"
	"booking-service/internal/database"
	"booking-service/internal/database/seeds"
	"booking-service/internal/middleware"
	"booking-service/pkg/utils"

	"booking-service/internal/handlers"
	"booking-service/internal/messaging"
	"booking-service/internal/repositories"
	"booking-service/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

// @title SeatGuard Booking Service API
// @version 1.0
// @description API para reservas, asientos, checkout y tickets.
// @contact.name Lucas Cabral
// @contact.email lucassimple@hotmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:4000
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Formato: Bearer <token>
// @tag.name Events
// @description Operaciones para gestionar eventos
// @tag.name Seats
// @description Operaciones para gestionar asientos
// @tag.name Booking Order
// @description Operaciones para gestionar órdenes de reserva 
// @tag.name Checkout
// @description Operaciones para gestionar el proceso de checkout

type SendMessageReq struct {
	Message string `json:"message" binding:"required"`
}

func main() {
	runSeed := flag.Bool("seed", false, "Run database seeding")
	runMigrate := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	cfg := config.LoadConfig()
	db := database.InitDB(context.Background(), cfg)

	// Si pongo "-migrate", ejecuto las migraciones y salgo
	if *runMigrate {
		log.Println("Ejecutando Migraciones de Base de Datos...")
		if err := database.RunMigrations(db); err != nil {
			log.Fatal("Error en la migración:", err)
		}
		return
	}

	// // Si pongo "-seed", ejecuto la semilla
	if *runSeed {
		log.Println("Ejecutando Semilla (Reset & Seed)...")
		if err := seeds.ResetAndSeed(db); err != nil {
			log.Fatal("Error seeding database", err)
		}
		log.Println("Database seeded successfully!")
		return
	}

	defer func() {
		if err := database.CloseDB(db); err != nil {
			fmt.Print("Base de datos cerrada")
		} else {
			fmt.Print("Error al cerrar la base de datos")
		}
	}()

	// Events
	eventRepo := repositories.NewEventRepository(db)
	eventService := services.NewEventService(eventRepo)
	eventHandler := handlers.NewEventHandler(eventService)

	// Seats
	seatRepo := repositories.NewSeatRepository(db)
	seatService := services.NewSeatService(seatRepo, eventRepo)
	seatHandler := handlers.NewSeatHandler(seatService)

	// Booking Orders
	bookingOrderRepo := repositories.NewBookingOrderRepository(db)
	bookingOrderService := services.NewBookingOrderService(bookingOrderRepo, seatRepo, eventRepo)
	bookingOrderHandler := handlers.NewBookingOrderHandler(bookingOrderService)

	// Checkout
	checkoutRepo := repositories.NewCheckoutRepository(db)
	checkoutService := services.NewCheckoutService(checkoutRepo)
	checkoutHandler := handlers.NewCheckoutHandler(checkoutService)

	// Ticket
	ticketRepo := repositories.NewTicketRepository(db)
	ticketService := services.NewTicketService(ticketRepo, bookingOrderRepo, seatRepo, eventRepo)
	pdfService := services.NewPDFService()
	ticketHandler := handlers.NewTicketHandler(ticketService, pdfService, bookingOrderService, checkoutService)

	// Emails
	host := cfg.Smtp_Host
	port := cfg.Smtp_Port
	user := cfg.Smtp_User
	pass := cfg.Smtp_Pass
	from := cfg.Smtp_From
	workers := cfg.Workers
	workersInt, err := strconv.Atoi(workers)
	if err != nil {
		log.Fatalf("Invalid WORKERS: %v", err)
	}
	// ---
	emailRepo, err := repositories.NewEmailRepository(host, port, user, pass, from)
	if err != nil {
		log.Fatalf("❌ Failed to initialize email repository: %v", err)
	}
	emailService := services.NewEmailService(emailRepo, workersInt)
	emailHandler := handlers.NewEmailHandler(emailService)

	// Queue AWS SQS
	ctx := context.Background()
	envs := config.LoadConfig()

	sqsClient, err := messaging.NewSQSClient(ctx, envs.AWSRegion, envs.SQSQueueUrl)
	if err != nil {
		log.Fatal(err)
	}

	sqsHandler := handlers.NewSQSHandler(sqsClient)

	guardUserJWT := middleware.UserMiddleware()

	globalUrl := middleware.NewRateLimiter(rate.Every(time.Second/10), 20)

	r := gin.Default()
	r.Use(utils.GetCorsConfig())
	r.Use(globalUrl.Middleware())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Health is OK!",
		})
	})
	v1 := r.Group("/api/v1")
	{
		events := v1.Group("/events")
		{
			// Events
			events.POST("", guardUserJWT, eventHandler.CreateEvent)
			events.GET("", guardUserJWT, eventHandler.GetAllEvents)
			events.GET("/:id", guardUserJWT, eventHandler.GetEventByID)
			events.PATCH("/:id", guardUserJWT, eventHandler.UpdateEvent)
			// Actualizo la disponibilidad de asientos de un evento. Se debe ejecutar con la confirmacion de un pago satisfactorio.
			events.PATCH("/availability/:id", eventHandler.UpdateAvailabilityForEvent)
			events.DELETE("/:id", guardUserJWT, eventHandler.DeleteEvent)
		}
		seats := v1.Group("/seats")
		{
			// Seats
			seats.POST("", guardUserJWT, seatHandler.CreateSeat)
			seats.GET("", seatHandler.GetSeats)
			seats.GET("/:id", guardUserJWT, seatHandler.GetSeat)
			seats.GET("/event/:eventId", seatHandler.GetSeatsByEventId)
			seats.PATCH("/:id", guardUserJWT, seatHandler.UpdateSeat)             // Para cambiar el estatus del asiento
			seats.PATCH("/lock/:id/uid/:uid", guardUserJWT, seatHandler.LockSeat) // Para bloquear un asiento
		}
		bookingOrders := v1.Group("/booking-orders")
		{
			// Booking Orders
			bookingOrders.POST("", guardUserJWT, bookingOrderHandler.CreateBookingOrder)
			bookingOrders.GET("", guardUserJWT, bookingOrderHandler.GetBookingOrders)
			bookingOrders.GET("/:id", guardUserJWT, bookingOrderHandler.GetBookingOrderById)
			bookingOrders.GET("/user/:id", guardUserJWT, bookingOrderHandler.GetAllOrderForUserID)
			bookingOrders.PATCH("/:id", guardUserJWT, bookingOrderHandler.UpdateBookingOrder)
		}
		checkouts := v1.Group("/checkouts")
		{
			checkouts.POST("", guardUserJWT, checkoutHandler.Create)
			checkouts.GET("/:orderID", guardUserJWT, checkoutHandler.GetByOrderID)
			checkouts.GET("", guardUserJWT, checkoutHandler.GetAll)
			checkouts.PUT("/:id", guardUserJWT, checkoutHandler.Update)
		}
		sqsMessaging := v1.Group("/sqs")
		{
			sqsMessaging.POST("/messaging", sqsHandler.Send)
		}
		// Creacion de checkout session
		stripe := v1.Group("/stripe")
		{
			stripe.POST("/create/checkout/session", guardUserJWT, handlers.CreateCartCheckoutSession(seatService, bookingOrderService))
		}
		// ✅ Generacion de ticket (NUEVO)
		tickets := v1.Group("/tickets")
		{
			// Crear Ticket para PDF
			tickets.POST("/", guardUserJWT, ticketHandler.CreateTicketFromEndpoint)
			// Obtener metadata del ticket
			tickets.GET("/:orderID", guardUserJWT, ticketHandler.GetTicketMetadata)
			// Descargar PDF del ticket (sin Bearer token)
			tickets.GET("/:orderID/download", ticketHandler.DownloadTicketPDF)
			// Regenerar PDF
			tickets.POST("/:orderID/regenerate", guardUserJWT, ticketHandler.RegenerateTicketPDF)
			// Admin: obtener todos los tickets
			tickets.GET("", guardUserJWT, ticketHandler.GetAllTickets)
			// Admin: obtener ticket por ID
			tickets.GET("/by-id/:ticketID", guardUserJWT, ticketHandler.GetTicketByID)
			// Admin: eliminar ticket
			tickets.DELETE("/:orderID", guardUserJWT, ticketHandler.DeleteTicket)
		}
		emails := v1.Group("/emails")
		{
			// Mandar un solo email (instantaneamente)
			emails.POST("/send", guardUserJWT, emailHandler.SendSync)
			// Mandar varios emails (masivo)
			emails.POST("/send-bulk", guardUserJWT, emailHandler.SendBulk)
			// Mandar varios emails en segundo plano (asíncrono, en cola)
			emails.POST("/send-bulk-async", guardUserJWT, emailHandler.SendAsync)
		}
	}

	log.Printf("Server starting on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
