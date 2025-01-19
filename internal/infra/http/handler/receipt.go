package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/john-moris/receipt-processor-challenge/internal/domain/model"
	"github.com/john-moris/receipt-processor-challenge/internal/domain/repository"
	"github.com/john-moris/receipt-processor-challenge/internal/infra/http/request"
	"github.com/john-moris/receipt-processor-challenge/internal/infra/http/response"
	"github.com/john-moris/receipt-processor-challenge/internal/infra/process"
	"github.com/labstack/echo/v4"
)

type Receipt struct {
	process    *process.Processor
	repository repository.Repository
}

func NewReceipt(r repository.Repository, p *process.Processor) *Receipt {
	return &Receipt{
		process:    p,
		repository: r,
	}
}

func (r *Receipt) Process(c echo.Context) error {
	var req request.Receipt

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	log.Println(req)

	var m model.Receipt

	m.Retailer = req.Retailer

	total, err := strconv.ParseFloat(req.Total, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	m.Total = total

	pt, err := time.Parse("2006-01-02 15:04", req.PurchaseDate+" "+req.PurchaseTime)
	if err != nil {
		return echo.ErrBadRequest
	}

	m.PurchaseTime = pt

	for _, item := range req.Items {
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return echo.ErrBadRequest
		}

		m.Items = append(m.Items, model.Item{
			Price:            price,
			ShortDescription: item.ShortDescription,
		})
	}

	log.Println(m)

	id := r.process.New(m)

	return c.JSON(http.StatusOK, response.ID{ID: id})
}

func (r *Receipt) Points(c echo.Context) error {
	id := c.Param("id")

	score, err := r.repository.Get(id)
	if err != nil {
		if errors.Is(err, repository.ErrItemNotFound) {
			return echo.ErrNotFound
		}

		if errors.Is(err, repository.ErrItemStillInProgress) {
			return c.NoContent(http.StatusAccepted)
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, response.Points{Points: score})
}

func (r *Receipt) Register(g *echo.Group) {
	g.POST("/process", r.Process)
	g.GET("/:id/points", r.Points)
}
