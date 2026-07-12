package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	)

	templates := template.Must(
		template.ParseFiles(
			"templates/layouts/base.html",
			"templates/pages/home.html",
		),
	)

	port := getEnv("PORT", "8080")
	addr := ":" + port

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title       string
			Description string
		}{
			Title:       "DGSIS — Engineering software that lasts.",
			Description: "DGSIS builds reliable digital systems designed to evolve, scale and remain maintainable for years.",
		}

		err := templates.ExecuteTemplate(
			w,
			"base",
			data,
		)

		if err != nil {

			logger.Error(
				"template rendering failed",
				"error",
				err,
			)

			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError,
			)

			return
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	logger.Info(
		"server starting",
		"address",
		addr,
	)

	server := &http.Server{
		Addr: addr,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error(
			"server stopped",
			"error",
			err,
		)

		os.Exit(1)
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
