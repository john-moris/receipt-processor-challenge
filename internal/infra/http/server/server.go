package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/john-moris/receipt-processor-challenge/internal/domain/repository"
	"github.com/john-moris/receipt-processor-challenge/internal/infra/http/handler"
	"github.com/john-moris/receipt-processor-challenge/internal/infra/process"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(lc fx.Lifecycle, logger *zap.Logger, r repository.Repository, p *process.Processor) *echo.Echo {
	e := echo.New()

	handler.NewReceipt(r, p).Register(e.Group("/receipts"))

	lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					if err := e.Start(":1378"); !errors.Is(err, http.ErrServerClosed) {
						logger.Fatal("echo initiation failed", zap.Error(err))
					}
				}()

				return nil
			},
			OnStop: e.Shutdown,
		},
	)

	return e
}
