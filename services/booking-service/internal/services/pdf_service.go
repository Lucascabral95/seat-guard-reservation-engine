package services

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"booking-service/internal/models"

	"github.com/jung-kurt/gofpdf"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) tr(text string) string {
	r := strings.NewReplacer(
		"á", "\xE1", "é", "\xE9", "í", "\xED", "ó", "\xF3", "ú", "\xFA",
		"ñ", "\xF1", "Ñ", "\xD1",
		"Á", "\xC1", "É", "\xC9", "Í", "\xCD", "Ó", "\xD3", "Ú", "\xDA",
		"ü", "\xFC", "Ü", "\xDC",
	)
	return r.Replace(text)
}

func formatMoney(amount int64) string {
	s := fmt.Sprintf("%d", amount)
	n := len(s)
	if n <= 3 {
		return s
	}
	out := ""
	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			out += "."
		}
		out += string(c)
	}
	return out
}

func (s *PDFService) GenerateTicket(ticket *models.TicketPDF) ([]byte, error) {
	if ticket == nil {
		return nil, errors.New("ticket cannot be nil")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 0) // ✅ GARANTIZA 1 SOLA PÁGINA
	pdf.SetTitle(s.tr("SeatGuards Ticket - "+ticket.EventName), false)
	pdf.SetAuthor("SeatGuards", false)
	pdf.AddPage()

	primary := []int{20, 24, 40}
	accent := []int{0, 170, 120}
	soft := []int{245, 247, 250}
	text := []int{33, 33, 33}
	muted := []int{120, 120, 120}

	// ================= HEADER =================
	pdf.SetFillColor(primary[0], primary[1], primary[2])
	pdf.Rect(0, 0, 210, 60, "F")

	pdf.SetFont("Helvetica", "B", 30)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(12, 12)
	pdf.Cell(0, 12, "SeatGuards")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(accent[0], accent[1], accent[2])
	pdf.SetXY(12, 26)
	pdf.Cell(0, 6, "Verified Digital Ticket")

	pdf.SetXY(12, 34)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetTextColor(180, 180, 180)
	pdf.Cell(20, 6, "ORDER ID:")

	pdf.SetFont("Courier", "", 10)
	pdf.SetTextColor(255, 255, 255)
	pdf.Cell(0, 6, ticket.OrderID)

	// ================= EVENTO =================
	pdf.SetY(75)
	pdf.SetFont("Helvetica", "B", 22)
	pdf.SetTextColor(text[0], text[1], text[2])
	pdf.Cell(0, 10, s.tr(ticket.EventName))

	pdf.SetDrawColor(accent[0], accent[1], accent[2])
	pdf.SetLineWidth(1.5)
	pdf.Line(12, 88, 198, 88)

	pdf.SetY(95)
	pdf.SetFont("Helvetica", "B", 10)

	pdf.SetTextColor(muted[0], muted[1], muted[2])
	pdf.Cell(35, 6, "FECHA")
	pdf.SetTextColor(text[0], text[1], text[2])
	pdf.Cell(65, 6, s.tr(time.Now().Format("02 Jan 2006")))

	pdf.SetTextColor(muted[0], muted[1], muted[2])
	pdf.Cell(20, 6, "HORA")
	pdf.SetTextColor(text[0], text[1], text[2])
	pdf.Cell(0, 6, ticket.EventHour)

	pdf.Ln(10)
	pdf.SetTextColor(muted[0], muted[1], muted[2])
	pdf.Cell(35, 6, s.tr("UBICACIÓN"))
	pdf.SetTextColor(text[0], text[1], text[2])
	pdf.Cell(0, 6, "SeatGuards Arena")

	// ================= COMPRADOR =================
	pdf.SetY(120)
	pdf.SetFillColor(soft[0], soft[1], soft[2])
	pdf.Rect(12, 120, 186, 26, "F")

	pdf.SetXY(16, 124)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetTextColor(muted[0], muted[1], muted[2])
	pdf.Cell(30, 5, "TITULAR")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(text[0], text[1], text[2])
	pdf.Cell(0, 5, s.tr(ticket.Name))

	pdf.SetXY(16, 132)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetTextColor(muted[0], muted[1], muted[2])
	pdf.Cell(30, 5, "EMAIL")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(text[0], text[1], text[2])
	pdf.Cell(0, 5, ticket.Email)

	// ================= TOTAL =================
	total := ticket.Amount / 100
	pdf.SetY(155)

	pdf.SetFillColor(primary[0], primary[1], primary[2])
	pdf.Rect(108, 155, 90, 16, "F")

	pdf.SetXY(112, 160)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(255, 255, 255)
	pdf.Cell(30, 6, "TOTAL")

	pdf.SetFont("Helvetica", "B", 15)
	pdf.SetTextColor(accent[0], accent[1], accent[2])
	pdf.Cell(0, 6, fmt.Sprintf("$%s %s", formatMoney(total), ticket.Currency))

	// ================= FOOTER =================
	pdf.SetY(280)
	pdf.SetFont("Helvetica", "", 7)
	pdf.SetTextColor(150, 150, 150)

	footer := fmt.Sprintf("Generado: %s | Version %d",
		time.Now().Format("2006-01-02 15:04:05"), ticket.PDFVersion)

	pdf.Cell(0, 4, s.tr(footer))

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
