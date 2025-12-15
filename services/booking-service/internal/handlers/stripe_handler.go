package handlers

import (
	"booking-service/internal/models"
	"booking-service/internal/services"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	checkoutsession "github.com/stripe/stripe-go/v76/checkout/session"
)

type SeatStruc struct {
	Id string `json:"id"`
}

type TicketItem struct {
	Name    string    `json:"name"`
	Amount  int64     `json:"amount"`
	SeatIds SeatStruc `json:"seatIds"`
}

type CreateCartCheckoutReq struct {
	UserId   string       `json:"userId"`
	Currency string       `json:"currency"`
	Items    []TicketItem `json:"items"`
}

// Inyectamos servicio de Seats
type PaymentHandler struct {
	seatService *services.SeatService
}

func NewPaymentHandlers(s *services.SeatService) *PaymentHandler {
	return &PaymentHandler{seatService: s}
}

// Creamos el checkout session, y si alguna de las entradas está vendida (SOLD), se denegará la compra. Esto evita duplicados.
func CreateCartCheckoutSession(seatService *services.SeatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

		var body CreateCartCheckoutReq
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body: " + err.Error()})
			return
		}

		var eventID string

		for _, item := range body.Items {
			seatID := item.SeatIds.Id
			if seatID != "" {
				// 1. Consultamos al servicio (ahora tenemos acceso a seatService)
				seat, err := seatService.GetSeat(seatID)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found: " + seatID})
					return
				}

				if eventID == "" {
					eventID = seat.EventID
				}

				// 2. Validamos estado
				if seat.Status != models.StatusAvailable {
					c.JSON(http.StatusConflict, gin.H{"error": "Seat/s is not available", "seatId": seatID})
					return
				}
			}
		}

		if len(body.Items) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
			return
		}

		currency := strings.ToLower(strings.TrimSpace(body.Currency))
		if currency == "" {
			currency = "usd"
		}

		var lineItems []*stripe.CheckoutSessionLineItemParams
		var allSeatIds []string

		for _, item := range body.Items {
			lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String(currency),
					UnitAmount: stripe.Int64(item.Amount * 100),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(item.Name),
					},
				},
				Quantity: stripe.Int64(1),
			})

			if item.SeatIds.Id != "" {
				allSeatIds = append(allSeatIds, item.SeatIds.Id)
			}
		}

		seatsMetadata := strings.Join(allSeatIds, ",")
		if len(seatsMetadata) > 450 {
			seatsMetadata = "many_seats_check_db"
		}

		params := &stripe.CheckoutSessionParams{
			Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
			LineItems:  lineItems,
			SuccessURL: stripe.String(os.Getenv("STRIPE_SUCCESS_URL")),
			CancelURL:  stripe.String(os.Getenv("STRIPE_CANCEL_URL")),
			Metadata: map[string]string{
				"user_id":  body.UserId,
				"seat_ids": seatsMetadata,
				"event_id": eventID,
			},
		}

		s, err := checkoutsession.New(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"url": s.URL})
	}
}
