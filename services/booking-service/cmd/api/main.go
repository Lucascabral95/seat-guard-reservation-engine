package main

import (
	"context"
	"flag"
	"log"

	"booking-service/internal/config"
	"booking-service/internal/database"
	"booking-service/internal/database/seeds"

	// "booking-service/internal/database/seeds"
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

	// Si pongo "-seed", corro la semilla
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
	seatService := services.NewSeatService(seatRepo)
	seatHandler := handlers.NewSeatHandler(seatService)

	// Booking Orders
	bookingOrderRepo := repositories.NewBookingOrderRepository(db)
	bookingOrderService := services.NewBookingOrderService(bookingOrderRepo)
	bookingOrderHandler := handlers.NewBookingOrderHandler(bookingOrderService)

	// Queue AWS SQS
	ctx := context.Background()
	envs := config.LoadConfig()

	sqsClient, err := messaging.NewSQSClient(ctx, envs.AWSRegion, envs.SQSQueueUrl)
	if err != nil {
		log.Fatal(err)
	}

	sqsHandler := handlers.NewSQSHandler(sqsClient)

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		events := v1.Group("/events")
		{
			// Events
			events.POST("", eventHandler.CreateEvent)
			events.GET("", eventHandler.GetAllEvents)
			events.GET("/:id", eventHandler.GetEventByID)
			events.PATCH("/:id", eventHandler.UpdateEvent)
			events.DELETE("/:id", eventHandler.DeleteEvent)
		}
		seats := v1.Group("/seats")
		{
			// Seats
			seats.POST("", seatHandler.CreateSeat)
			seats.GET("", seatHandler.GetSeats)
			seats.GET("/:id", seatHandler.GetSeat)
			seats.GET("/event/:eventId", seatHandler.GetSeatsByEventId)
			seats.PATCH("/:id", seatHandler.UpdateSeat)
			seats.PATCH("/lock/:id/uid/:uid", seatHandler.LockSeat)
		}
		bookingOrders := v1.Group("/booking-orders")
		{
			// Booking Orders
			bookingOrders.POST("", bookingOrderHandler.CreateBookingOrder)
			bookingOrders.GET("", bookingOrderHandler.GetBookingOrders)
			bookingOrders.GET("/:id", bookingOrderHandler.GetBookingOrderById)
		}
		sqsMessaging := v1.Group("/sqs")
		{
			sqsMessaging.POST("/messaging", sqsHandler.Send)
		}
		// Creacion de checkout session
		stripe := v1.Group("/stripe")
		{
			stripe.POST("/create/checkout/session", handlers.CreateCartCheckoutSession)
		}
	}

	log.Printf("Server starting on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
