package main

import (
	"computerextra/datenschutz_training_golang/internal/app"
	"context"
	"embed"
	"log/slog"
	"os"
	"os/signal"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var (
	files embed.FS
	store = sessions.NewCookieStore([]byte("super-secret-key"))
)

func main() {
	godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app, err := app.New(logger, app.Config{}, files, store)
	if err != nil {
		logger.Error("failed to create app", slog.Any("error", err))
	}

	if err := app.Start(ctx); err != nil {
		logger.Error("failed to start app", slog.Any("error", err))
	}
}
