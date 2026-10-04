package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

type pageData struct {
	Title       string
	Description string
	CurrentPath string
}

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	)

	homeTmpl := mustParsePage("templates/pages/home.html")
	engineeringTmpl := mustParsePage("templates/pages/engineering.html")
	workTmpl := mustParsePage("templates/pages/work.html")
	contactTmpl := mustParsePage("templates/pages/contact.html")

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
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		renderPage(w, r, logger, homeTmpl, pageData{
			Title:       "DGSIS — Engineering software that lasts.",
			Description: "DGSIS is a software engineering company that builds reliable digital systems designed to evolve, scale and remain maintainable for years.",
			CurrentPath: "/",
		})
	})

	http.HandleFunc("/engineering", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, r, logger, engineeringTmpl, pageData{
			Title:       "Engineering — DGSIS",
			Description: "DGSIS designs and builds backend systems, software architecture and digital platforms focused on reliability and long-term maintainability.",
			CurrentPath: "/engineering",
		})
	})

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, r, logger, workTmpl, pageData{
			Title:       "Work — DGSIS",
			Description: "Selected engineering work from DGSIS. Case studies focused on problems, decisions, solutions and lasting results.",
			CurrentPath: "/work",
		})
	})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, r, logger, contactTmpl, pageData{
			Title:       "Contact — DGSIS",
			Description: "Start a conversation with DGSIS about building digital systems with solid engineering foundations.",
			CurrentPath: "/contact",
		})
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

func mustParsePage(pagePath string) *template.Template {
	return template.Must(
		template.ParseFiles(
			"templates/layouts/base.html",
			pagePath,
		),
	)
}

func renderPage(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	tmpl *template.Template,
	data pageData,
) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := tmpl.ExecuteTemplate(
		w,
		"base",
		data,
	)

	if err != nil {
		logger.Error(
			"template rendering failed",
			"error",
			err,
			"path",
			r.URL.Path,
		)

		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
