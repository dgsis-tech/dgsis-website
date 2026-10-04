(function () {
	"use strict";

	var THEME_KEY = "dgsis_theme";
	var LANG_KEY = "dgsis_lang";
	var THEMES = { light: true, dark: true, system: true };
	var THEME_COLORS = {
		light: "#F3F7FB",
		dark: "#081726",
	};

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

	function syncSummaryIcon(theme) {
		var root = document.documentElement;
		var icons = root.querySelectorAll("[data-prefs-theme-icon]");
		for (var i = 0; i < icons.length; i++) {
			icons[i].setAttribute("data-active-theme", theme);
		}
	}

	function applyTheme(theme) {
		var root = document.documentElement;
		root.setAttribute("data-theme", theme);

		var scheme = resolvedScheme(theme);
		var colorMeta = document.getElementById("meta-theme-color");
		if (colorMeta) {
			colorMeta.setAttribute("content", THEME_COLORS[scheme] || THEME_COLORS.dark);
		}

		var buttons = document.querySelectorAll("[data-theme-value]");
		for (var i = 0; i < buttons.length; i++) {
			var btn = buttons[i];
			var active = btn.getAttribute("data-theme-value") === theme;
			btn.classList.toggle("is-active", active);
			btn.setAttribute("aria-pressed", active ? "true" : "false");
		}

		syncSummaryIcon(theme);
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

	function closePrefsMenus(except) {
		var menus = document.querySelectorAll("details.prefs-menu");
		for (var i = 0; i < menus.length; i++) {
			if (menus[i] !== except) {
				menus[i].open = false;
			}
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

		document.addEventListener("click", function (event) {
			var menu = event.target.closest("details.prefs-menu");
			if (!menu) {
				closePrefsMenus(null);
				return;
			}
			closePrefsMenus(menu);
		});

		document.addEventListener("keydown", function (event) {
			if (event.key === "Escape") {
				closePrefsMenus(null);
			}
		});

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
