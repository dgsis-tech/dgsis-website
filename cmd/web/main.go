package main

import (
	"compress/gzip"
	"fmt"
	"html"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type pageData struct {
	Title        string
	Description  string
	CurrentPath  string
	FormFeedback template.HTML
	Lang         string
	Msg          Msg
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

	mux := http.NewServeMux()

	mux.Handle(
		"/static/",
		cacheControl(
			"public, max-age=86400",
			http.StripPrefix(
				"/static/",
				http.FileServer(http.Dir("static")),
			),
		),
	)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		lang, redirected := resolveLang(w, r)
		if redirected {
			return
		}
		msg := messages(lang)

		renderPage(w, r, logger, homeTmpl, pageData{
			Title:       msg.HomeTitle,
			Description: msg.HomeDescription,
			CurrentPath: "/",
			Lang:        lang,
			Msg:         msg,
		})
	})

	mux.HandleFunc("/engineering", func(w http.ResponseWriter, r *http.Request) {
		lang, redirected := resolveLang(w, r)
		if redirected {
			return
		}
		msg := messages(lang)

		renderPage(w, r, logger, engineeringTmpl, pageData{
			Title:       msg.EngTitle,
			Description: msg.EngDescription,
			CurrentPath: "/engineering",
			Lang:        lang,
			Msg:         msg,
		})
	})

	mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		lang, redirected := resolveLang(w, r)
		if redirected {
			return
		}
		msg := messages(lang)

		renderPage(w, r, logger, workTmpl, pageData{
			Title:       msg.WorkTitle,
			Description: msg.WorkDescription,
			CurrentPath: "/work",
			Lang:        lang,
			Msg:         msg,
		})
	})

	mux.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
			lang, redirected := resolveLang(w, r)
			if redirected {
				return
			}
			msg := messages(lang)

			feedback := template.HTML("")
			if r.URL.Query().Get("sent") == "1" {
				feedback = template.HTML(
					`<p class="contact-form__ok">` +
						html.EscapeString(msg.ContactOK) +
						`</p>`,
				)
			}

			renderPage(w, r, logger, contactTmpl, pageData{
				Title:        msg.ContactTitle,
				Description:  msg.ContactDescription,
				CurrentPath:  "/contact",
				FormFeedback: feedback,
				Lang:         lang,
				Msg:          msg,
			})

		case http.MethodPost:
			handleContactSubmit(w, r, logger, contactTmpl)

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")

		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	logger.Info(
		"server starting",
		"address",
		addr,
	)

	server := &http.Server{
		Addr:              addr,
		Handler:           withSecurityHeaders(withGzip(mux)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Language", data.Lang)

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

func handleContactSubmit(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	contactTmpl *template.Template,
) {
	// POST must not redirect for ?lang=; use cookie / Accept-Language / default ES.
	lang := langFromRequest(r)
	msg := messages(lang)

	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)

	if err := r.ParseForm(); err != nil {
		writeContactFeedback(w, r, logger, contactTmpl, lang, msg, false, msg.ContactErrRead)
		return
	}

	// Honeypot
	if strings.TrimSpace(r.FormValue("website")) != "" {
		writeContactFeedback(w, r, logger, contactTmpl, lang, msg, true, msg.ContactOKShort)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	company := strings.TrimSpace(r.FormValue("company"))
	message := strings.TrimSpace(r.FormValue("message"))

	if utf8.RuneCountInString(name) < 2 || utf8.RuneCountInString(name) > 80 {
		writeContactFeedback(w, r, logger, contactTmpl, lang, msg, false, msg.ContactErrName)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil || utf8.RuneCountInString(email) > 120 {
		writeContactFeedback(w, r, logger, contactTmpl, lang, msg, false, msg.ContactErrEmail)
		return
	}

	if utf8.RuneCountInString(company) > 120 {
		writeContactFeedback(w, r, logger, contactTmpl, lang, msg, false, msg.ContactErrCompany)
		return
	}

	if utf8.RuneCountInString(message) < 20 || utf8.RuneCountInString(message) > 4000 {
		writeContactFeedback(w, r, logger, contactTmpl, lang, msg, false, msg.ContactErrMessage)
		return
	}

	logger.Info(
		"contact inquiry received",
		"name", name,
		"email", email,
		"company", company,
		"message_length", utf8.RuneCountInString(message),
		"lang", lang,
	)

	writeContactFeedback(
		w,
		r,
		logger,
		contactTmpl,
		lang,
		msg,
		true,
		msg.ContactOK,
	)
}

func langFromRequest(r *http.Request) string {
	if c, err := r.Cookie(langCookieName); err == nil {
		if v := normalizeLang(c.Value); v != "" {
			return v
		}
	}
	for _, part := range strings.Split(r.Header.Get("Accept-Language"), ",") {
		token := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		if v := normalizeLang(token); v != "" {
			return v
		}
	}
	return langES
}

func writeContactFeedback(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	contactTmpl *template.Template,
	lang string,
	msg Msg,
	ok bool,
	message string,
) {
	class := "contact-form__err"
	if ok {
		class = "contact-form__ok"
	}

	fragment := fmt.Sprintf(
		`<p class="%s">%s</p>`,
		class,
		html.EscapeString(message),
	)

	if r.Header.Get("HX-Request") == "true" {
		if !ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, fragment)
		return
	}

	if ok {
		http.Redirect(w, r, "/contact?sent=1", http.StatusSeeOther)
		return
	}

	renderPage(w, r, logger, contactTmpl, pageData{
		Title:        msg.ContactTitle,
		Description:  msg.ContactDescription,
		CurrentPath:  "/contact",
		FormFeedback: template.HTML(fragment),
		Lang:         lang,
		Msg:          msg,
	})
}

func cacheControl(value string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", value)
		next.ServeHTTP(w, r)
	})
}

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func withGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		next.ServeHTTP(gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
