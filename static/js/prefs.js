(function () {
	"use strict";

	var THEME_KEY = "dgsis_theme";
	var LANG_KEY = "dgsis_lang";
	var THEMES = { light: true, dark: true, system: true };

	function readTheme() {
		try {
			var stored = localStorage.getItem(THEME_KEY);
			if (stored && THEMES[stored]) {
				return stored;
			}
		} catch (e) {
			/* private mode / blocked storage */
		}
		return "system";
	}

	function resolvedScheme(theme) {
		if (theme === "light" || theme === "dark") {
			return theme;
		}
		if (window.matchMedia && window.matchMedia("(prefers-color-scheme: light)").matches) {
			return "light";
		}
		return "dark";
	}

	function applyTheme(theme) {
		var root = document.documentElement;
		root.setAttribute("data-theme", theme);

		var scheme = resolvedScheme(theme);
		var colorMeta = document.getElementById("meta-theme-color");
		if (colorMeta) {
			colorMeta.setAttribute("content", scheme === "light" ? "#eef1f5" : "#0b0d10");
		}

		var buttons = document.querySelectorAll("[data-theme-value]");
		for (var i = 0; i < buttons.length; i++) {
			var btn = buttons[i];
			var active = btn.getAttribute("data-theme-value") === theme;
			btn.classList.toggle("is-active", active);
			btn.setAttribute("aria-pressed", active ? "true" : "false");
		}
	}

	function persistTheme(theme) {
		try {
			localStorage.setItem(THEME_KEY, theme);
		} catch (e) {
			/* ignore */
		}
		try {
			document.cookie =
				THEME_KEY +
				"=" +
				encodeURIComponent(theme) +
				"; path=/; max-age=31536000; samesite=lax";
		} catch (e2) {
			/* ignore */
		}
	}

	function syncLangStorage() {
		var lang = document.documentElement.lang;
		if (lang !== "es" && lang !== "en") {
			return;
		}
		try {
			localStorage.setItem(LANG_KEY, lang);
		} catch (e) {
			/* ignore */
		}
	}

	// Apply before paint when possible (script is in <head> without defer).
	applyTheme(readTheme());

	document.addEventListener("DOMContentLoaded", function () {
		applyTheme(readTheme());
		syncLangStorage();

		var switcher = document.querySelector(".theme-switch");
		if (switcher) {
			switcher.addEventListener("click", function (event) {
				var target = event.target.closest("[data-theme-value]");
				if (!target) {
					return;
				}
				var value = target.getAttribute("data-theme-value");
				if (!THEMES[value]) {
					return;
				}
				persistTheme(value);
				applyTheme(value);
			});
		}

		if (window.matchMedia) {
			var mq = window.matchMedia("(prefers-color-scheme: dark)");
			var onChange = function () {
				if (readTheme() === "system") {
					applyTheme("system");
				}
			};
			if (mq.addEventListener) {
				mq.addEventListener("change", onChange);
			} else if (mq.addListener) {
				mq.addListener(onChange);
			}
		}
	});
})();
