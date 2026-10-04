package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNormalizeLang(t *testing.T) {
	cases := map[string]string{
		"es":      langES,
		"ES":      langES,
		"es-ES":   langES,
		"en":      langEN,
		"en-US":   langEN,
		" fr ":    "",
		"":        "",
		"espanol": "",
	}
	for in, want := range cases {
		if got := normalizeLang(in); got != want {
			t.Fatalf("normalizeLang(%q)=%q want %q", in, got, want)
		}
	}
}

func TestResolveLangCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: langCookieName, Value: "en"})
	w := httptest.NewRecorder()

	lang, redirected := resolveLang(w, r)
	if redirected {
		t.Fatal("unexpected redirect")
	}
	if lang != langEN {
		t.Fatalf("got %q want en", lang)
	}
}

func TestResolveLangAcceptLanguage(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-Language", "fr-FR,en;q=0.8,es;q=0.7")
	w := httptest.NewRecorder()

	lang, redirected := resolveLang(w, r)
	if redirected {
		t.Fatal("unexpected redirect")
	}
	if lang != langEN {
		t.Fatalf("got %q want en from Accept-Language", lang)
	}
}

func TestResolveLangDefaultES(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	lang, redirected := resolveLang(w, r)
	if redirected {
		t.Fatal("unexpected redirect")
	}
	if lang != langES {
		t.Fatalf("got %q want es default", lang)
	}
}

func TestResolveLangQuerySetsCookieAndRedirects(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/contact?lang=en&sent=1", nil)
	w := httptest.NewRecorder()

	lang, redirected := resolveLang(w, r)
	if !redirected {
		t.Fatal("expected redirect")
	}
	if lang != langEN {
		t.Fatalf("got %q want en", lang)
	}
	if loc := w.Header().Get("Location"); loc != "/contact?sent=1" {
		t.Fatalf("Location=%q", loc)
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != langCookieName || cookies[0].Value != langEN {
		t.Fatalf("cookie=%v", cookies)
	}
}

func TestMessagesBothLanguages(t *testing.T) {
	es := messages(langES)
	en := messages(langEN)
	if es.NavHome == en.NavHome {
		t.Fatal("expected different nav labels")
	}
	if es.HomeH1 == "" || en.HomeH1 == "" {
		t.Fatal("empty hero copy")
	}
}
