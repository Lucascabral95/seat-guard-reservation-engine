package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"booking-service/internal/models"
	"booking-service/internal/services"

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

type ResponseCartCheckoutReq struct {
	OrderBookingId string       `json:"orderBookingId"`
	UserId         string       `json:"userId"`
	Currency       string       `json:"currency"`
	Items          []TicketItem `json:"items"`
}

// CreateCartCheckoutSession Crea una sesión de pago en Stripe para un carrito de tickets
// @Summary Crear sesión de pago Stripe para carrito
// @Description Crea una sesión de pago en Stripe para un carrito de tickets
// @Tags Stripe
// @Accept json
// @Produce json
// @Param body body CreateCartCheckoutReq true "Datos del carrito"
// @Success 200 {object} map[string]string "Sesión de pago creada exitosamente"
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 404 {object} map[string]string "Asiento no encontrado"
// @Failure 409 {object} map[string]string "Asiento no disponible"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /stripe/cart/checkout [post]
// @Security BearerAuth
func CreateCartCheckoutSession(seatService *services.SeatService, orderService *services.BookingOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {

		stripeKey := os.Getenv("STRIPE_SECRET_KEY")
		if stripeKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "STRIPE_SECRET_KEY missing"})
			return
		}
		stripe.Key = stripeKey

		var body CreateCartCheckoutReq
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body: " + err.Error()})
			return
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
		var totalAmount int64 = 0
		var eventID string

		var enrichedItems []TicketItem

		for _, item := range body.Items {
			seatID := item.SeatIds.Id
			if seatID == "" {
				continue
			}

			seat, err := seatService.GetSeat(seatID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found: " + seatID})
				return
			}

			if eventID == "" {
				eventID = seat.EventID
			} else if eventID != seat.EventID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot mix events"})
				return
			}

			if seat.Status != models.StatusAvailable {
				c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Seat %s not available", seat.Number)})
				return
			}

			realAmount := int64(seat.Price * 100)
			realName := fmt.Sprintf("Asiento %s - %s", seat.Number, seat.Section)

			totalAmount += realAmount
			allSeatIds = append(allSeatIds, seatID)

			// Bloqueo
			seatService.LockSeat(seatID, body.UserId)

			enrichedItems = append(enrichedItems, TicketItem{
				Name:   realName,
				Amount: realAmount,
				SeatIds: SeatStruc{
					Id: seatID,
				},
			})

			lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String(currency),
					UnitAmount: stripe.Int64(realAmount),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(realName),
					},
				},
				Quantity: stripe.Int64(1),
			})
		}

		if len(allSeatIds) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No valid seats provided"})
			return
		}

		body.Items = enrichedItems

		order := &models.BookingOrder{
			UserID:  body.UserId,
			Status:  models.PaymentPending,
			SeatIDs: allSeatIds,
			Amount:  totalAmount,
		}

		if err := orderService.CreateBookingOrder(order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}

		seatsMetadata := strings.Join(allSeatIds, ",")
		if len(seatsMetadata) > 450 {
			seatsMetadata = "many_seats_check_db"
		}

		baseURL := os.Getenv("STRIPE_SUCCESS_URL")
		successURLWithParam := fmt.Sprintf("%s/dentro/checkout/success?session_id={CHECKOUT_SESSION_ID}&order_id=%s", baseURL, order.ID)
		errorURL := fmt.Sprintf("%s/dentro/checkout/cancel", baseURL)

		params := &stripe.CheckoutSessionParams{
			Mode:      stripe.String(string(stripe.CheckoutSessionModePayment)),
			LineItems: lineItems,

			SuccessURL: stripe.String(successURLWithParam),
			CancelURL:  stripe.String(errorURL),
			PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
				Metadata: map[string]string{
					"user_id":  body.UserId,
					"seat_ids": seatsMetadata,
					"event_id": eventID,
					"order_id": order.ID,
				},
			},
			Metadata: map[string]string{
				"user_id":  body.UserId,
				"seat_ids": seatsMetadata,
				"event_id": eventID,
				"order_id": order.ID,
			},
		}

		responsePayload := &ResponseCartCheckoutReq{
			OrderBookingId: order.ID,
			UserId:         body.UserId,
			Currency:       body.Currency,
			Items:          body.Items,
		}

		s, err := checkoutsession.New(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"url":          s.URL,
			"dataCheckout": BuildCheckoutResponse(responsePayload),
		})
	}
}

func BuildCheckoutResponse(callback *ResponseCartCheckoutReq) gin.H {
	return gin.H{
		"orderBookingId": callback.OrderBookingId,
		"userId":         callback.UserId,
		"currency":       callback.Currency,
		"items":          callback.Items,
	}
}
