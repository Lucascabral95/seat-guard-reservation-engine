package main

import (
	"context"
	"flag"
	"log"

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
)

type SendMessageReq struct {
	Message string `json:"message" binding:"required"`
}

func main() {
	runSeed := flag.Bool("seed", false, "Run database seeding")
	flag.Parse()

	cfg := config.LoadConfig()
	db := database.InitDB(cfg)

	// Si pongo "-seed", ejecuto la semilla
	if *runSeed {
		log.Println("Ejecutando Semilla (Reset & Seed)...")
		if err := seeds.ResetAndSeed(db); err != nil {
			log.Fatal("Error seeding database", err)
		}
		log.Println("Database seeded successfully!")
		return
	}

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

	// Queue AWS SQS
	ctx := context.Background()
	envs := config.LoadConfig()

	sqsClient, err := messaging.NewSQSClient(ctx, envs.AWSRegion, envs.SQSQueueUrl)
	if err != nil {
		log.Fatal(err)
	}

	sqsHandler := handlers.NewSQSHandler(sqsClient)

	guardUserJWT := middleware.UserMiddleware()

	r := gin.Default()
	r.Use(utils.GetCorsConfig())

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
	}

	log.Printf("Server starting on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
