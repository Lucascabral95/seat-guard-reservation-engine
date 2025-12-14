package handlers

import (
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

func CreateCartCheckoutSession(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

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

	for _, item := range body.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:   stripe.String(currency),
				UnitAmount: stripe.Int64(item.Amount),
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
		},
	}

	s, err := checkoutsession.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": s.URL})
}
