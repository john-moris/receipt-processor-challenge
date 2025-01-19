package main

import (
	"github.com/john-moris/receipt-processor-challenge/internal/domain/repository"
	"github.com/john-moris/receipt-processor-challenge/internal/infra/db"
	"github.com/john-moris/receipt-processor-challenge/internal/infra/http/server"
	"github.com/john-moris/receipt-processor-challenge/internal/infra/process"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(zap.NewDevelopment),
		fx.Provide(
			fx.Annotate(db.NewMemory, fx.As(new(repository.Repository))),
		),
		fx.Provide(process.New),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Provide(server.New),
		fx.Invoke(func(_ *echo.Echo) {
			_ = pterm.DefaultBigText.WithLetters(
				putils.LettersFromStringWithStyle("Fetch", pterm.FgLightCyan.ToStyle()),
				putils.LettersFromStringWithStyle(" Assignment", pterm.FgLightMagenta.ToStyle()),
				putils.LettersFromStringWithStyle(" 2025", pterm.FgLightRed.ToStyle()),
			).Render()
		}),
	).Run()
}
