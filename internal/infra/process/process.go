package process

import (
	"math"
	"strings"
	"unicode"

	"github.com/john-moris/receipt-processor-challenge/internal/domain/model"
	"github.com/john-moris/receipt-processor-challenge/internal/domain/repository"
	"go.uber.org/zap"
)

type Processor struct {
	r      repository.Repository
	logger *zap.Logger
}

func New(r repository.Repository, logger *zap.Logger) *Processor {
	return &Processor{
		r:      r,
		logger: logger,
	}
}

func (p *Processor) New(r model.Receipt) string {
	id := p.r.Start()

	go p.process(r, id)

	return id
}

func (p *Processor) process(r model.Receipt, id string) {
	score := 0

	for _, c := range r.Retailer {
		if unicode.IsLetter(c) {
			score++
		}
	}

	p.logger.Info("retailer name", zap.Int("score", score), zap.String("id", id))

	if math.Round(r.Total) == r.Total {
		score += 50
	}

	p.logger.Info("total is integer", zap.Int("score", score), zap.String("id", id))

	if math.Round(r.Total*4) == r.Total*4 { // nolint: mnd
		score += 25
	}

	p.logger.Info("total is dividable by 0.25", zap.Int("score", score), zap.String("id", id))

	score += (len(r.Items) / 2) * 5 // nolint: mnd
	p.logger.Info("items pairs", zap.Int("score", score), zap.String("id", id))

	if r.PurchaseTime.Day()%2 == 1 {
		score += 6
	}

	p.logger.Info("day is odd", zap.Int("score", score), zap.String("id", id))

	if r.PurchaseTime.Hour() >= 14 && r.PurchaseTime.Hour() < 16 {
		score += 10
	}

	p.logger.Info("hour is in 14 - 16", zap.Int("score", score), zap.String("id", id))

	for _, item := range r.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			score += int(math.Ceil(item.Price * 0.2)) // nolint: mnd
		}

		p.logger.Info("item description", zap.Int("score", score), zap.String("id", id))
	}

	p.logger.Info("processing finished", zap.String("id", id), zap.Int("score", score))

	p.r.Finish(id, score)
}
