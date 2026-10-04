# DGSIS Website — Technical Architecture Document v1

**Documento:** 04_TECHNICAL_ARCHITECTURE.md  
**Versión:** v1.0  
**Estado:** Fuente de verdad alineada con la implementación actual  
**Alcance:** Arquitectura real del repositorio `dgsis-website` (branch de producto / PR de contenido)

---

# 1. Architecture Overview

## Tipo de aplicación

DGSIS Website es una aplicación web **server-rendered (SSR)** escrita en **Go**, que genera HTML en el servidor mediante `html/template` y sirve CSS e imágenes estáticas.

No es una SPA. No hay framework frontend. No hay base de datos. No hay API de negocio más allá de un endpoint de salud.

## Objetivo arquitectónico

Demostrar, con la propia implementación, la filosofía DGSIS:

> Engineering software that lasts.

La arquitectura prioriza:

* claridad estructural
* dependencias mínimas
* mantenibilidad a largo plazo
* HTML como base de la experiencia
* operativa simple (binario + assets + Docker)

## Principios que guían las decisiones

1. **Standard library first** — usar la biblioteca estándar de Go mientras resuelva el problema.
2. **HTML first** — la experiencia nace del HTML semántico; la interactividad es opcional.
3. **No abstracción prematura** — no crear capas (`domain`, `services`, partials) sin necesidad real.
4. **CSS propio** — control total del sistema visual sin frameworks CSS.
5. **Evolución incremental** — crecer páginas y comportamiento solo cuando aporten valor.

## ¿Por qué representa la filosofía DGSIS?

Porque el sistema es deliberadamente simple, comprensible y operable. Las decisiones técnicas (Go + SSR + templates + CSS propio) tienen razón de ser: menos complejidad accidental, más claridad y una base que puede evolucionar sin reescribir el producto.

---

# 2. Technical Vision

## Qué demuestra esta web

La web demuestra ingeniería profesional mediante:

* una estructura fácil de leer
* un runtime predecible
* ausencia de moda tecnológica innecesaria
* coherencia entre mensaje de marca e implementación

## Principios de ingeniería representados

| Principio | Cómo se refleja |
| --- | --- |
| Mantenibilidad | Pocos archivos, responsabilidades claras, sin frameworks ocultos |
| Claridad | Handlers y templates legibles de extremo a extremo |
| Evolución incremental | Páginas nuevas se añaden como templates + rutas |
| Bajo acoplamiento | Layout global separado de páginas; assets estáticos independientes |
| Simplicidad operacional | Un binario, `PORT`, Docker multi-stage |

## Equilibrio simplicidad / profesionalidad

La simplicidad no implica improvisación. El proyecto incluye:

* metadata SEO básica
* landmarks y atributos de accesibilidad
* logging estructurado
* health check
* imagen Docker con usuario no root

La profesionalidad está en la disciplina, no en la cantidad de piezas.

---

# 3. Technology Stack

## Backend

| Elemento | Valor actual |
| --- | --- |
| Lenguaje | Go (`go.mod` → `go 1.26`) |
| Servidor HTTP | `net/http` (biblioteca estándar) |
| Templates | `html/template` |
| Logging | `log/slog` (JSON a stdout) |
| Configuración | Variable de entorno `PORT` (default `8080`) |
| Dependencias externas Go | Ninguna |

### ¿Por qué Go?

Go aporta simplicidad, un runtime claro, buen rendimiento y herramientas robustas. En este proyecto no se usa como “stack de moda”, sino como herramienta adecuada para un servidor HTTP pequeño, estable y fácil de operar.

## Frontend

| Elemento | Estado actual |
| --- | --- |
| HTML | Semántico, SSR |
| CSS | Propio en `static/css/style.css` |
| JavaScript | Mínimo: HTMX vendorizado en `static/js/htmx.min.js` |
| HTMX | Usado para el formulario de contacto (`POST /contact` → fragmento HTML) |

### ¿Por qué este enfoque?

La web comunica ingeniería. Un frontend mínimo:

* reduce superficie de fallo
* mejora comprensión
* mantiene el HTML como contrato principal
* evita convertir el sitio en una aplicación frontend compleja

## Infrastructure

| Elemento | Estado actual |
| --- | --- |
| Docker | Implementado (`Dockerfile` multi-stage, Alpine, usuario no root) |
| Nginx / Traefik | **No implementados** en el repositorio |
| CI/CD | **No implementado** en el repositorio |
| Orquestación | Fuera de alcance actual |

Despliegue actual: construir imagen y ejecutar el contenedor exponiendo el puerto `8080`.

---

# 4. Application Architecture

## Tipo de arquitectura

Arquitectura **monolítica mínima** centrada en un proceso HTTP:

```
Browser
  ↓
net/http ServeMux
  ↓
Handler (ruta)
  ↓
html/template render  ó  FileServer / JSON
  ↓
HTTP Response
```

## Separación de responsabilidades

| Pieza | Responsabilidad |
| --- | --- |
| `cmd/web/main.go` | Arranque, rutas, logging, render |
| `templates/layouts/base.html` | Documento HTML, header, footer, metadata |
| `templates/pages/*.html` | Contenido de cada página (`content`) |
| `static/` | CSS e imágenes públicas |
| `Dockerfile` | Empaquetado de producción |

## Flujo de una request de página

1. El cliente solicita una ruta (`/`, `/engineering`, `/work`, `/contact`).
2. El mux de `net/http` selecciona el handler.
3. El handler valida método (`GET`/`HEAD`).
4. Se ejecuta el template correspondiente (`base` + `content`).
5. Se devuelve HTML; errores de render se registran con `slog` y responden `500`.

## Flujo de assets

1. Request a `/static/...`
2. `http.FileServer` sirve desde el directorio `static/`
3. Sin procesamiento adicional

## Flujo de health

1. `GET /health`
2. Respuesta JSON `{"status":"ok"}`

---

# 5. Project Structure

Estructura real relevante:

```
.
├── cmd/web/main.go
├── docs/
│   ├── 01_BRAND_BOOK.md
│   ├── 02_UX_STRATEGY.md
│   ├── 03_DESIGN_SYSTEM.md
│   ├── 04_TECHNICAL_ARCHITECTURE.md
│   └── 05_IMPLEMENTATION_STATUS.md
├── static/
│   ├── css/style.css
│   └── images/favicon.svg
├── templates/
│   ├── layouts/base.html
│   └── pages/
│       ├── home.html
│       ├── engineering.html
│       ├── work.html
│       └── contact.html
├── Dockerfile
├── go.mod
└── README.md
```

## Responsabilidades

### `cmd/web`

Punto de entrada. Existe porque la aplicación es un binario Go convencional.

### `templates/layouts`

Layout global compartido. Existe para no duplicar `<html>`, navegación y footer.

### `templates/pages`

Una plantilla por página. Cada una define el bloque `content`.

### `static`

Assets versionados con el código. CSS e imágenes sin pipeline de bundling.

### `docs`

Fuente de verdad de marca, UX, diseño, arquitectura y estado de implementación.

### ¿Por qué no existe `internal/`?

No hay lógica de dominio, persistencia ni casos de uso que justifiquen packages internos. Introducirlos ahora sería abstracción prematura.

---

# 6. Go Application Design

## Organización

Todo el runtime vive en `package main` dentro de `cmd/web/main.go`.

## Piezas principales

### Configuración

```go
port := getEnv("PORT", "8080")
```

Única configuración operativa actual.

### Templates

Cada página se parsea en un set independiente:

```go
mustParsePage("templates/pages/home.html")
```

Esto evita colisiones del nombre `content` entre páginas.

### Datos de página

```go
type pageData struct {
	Title       string
	Description string
	CurrentPath string
}
```

`CurrentPath` alimenta `aria-current` en la navegación.

### Handlers

Handlers por ruta registrados con `http.HandleFunc`:

* `/` (exacta)
* `/engineering`
* `/work`
* `/contact`
* `/health`
* `/static/` (FileServer)

### Servidor

```go
server := &http.Server{Addr: addr}
server.ListenAndServe()
```

### Logging

`slog` con handler JSON a stdout para observabilidad básica en contenedor/proceso.

### Convenciones

* Nombres claros y funciones pequeñas (`mustParsePage`, `renderPage`, `getEnv`)
* Sin interfaces inventadas
* Sin inyección de dependencias formal (no necesaria aún)

---

# 7. Template Architecture

## Estrategia

```
layouts/base.html  → define "base"
pages/*.html       → definen "content"
base incluye         {{template "content" .}}
```

## Responsabilidades

| Template | Responsabilidad |
| --- | --- |
| `base.html` | HTML shell, metadata, nav, main, footer |
| `home.html` | Narrativa completa Home (Hero → Approach → Engineering → Work → Contact) |
| `engineering.html` | Profundización de capacidades y método |
| `work.html` | Estructura de case studies (sin inventar proyectos) |
| `contact.html` | Conversación / CTA mailto |

## Cuándo extraer sections o partials

Solo cuando exista reutilización real y repetida. Hoy la duplicación entre páginas es baja y deliberada.

## Regla

No abstraer antes de necesitarlo.

---

# 8. HTMX Strategy

## Estado

HTMX **está integrado** de forma puntual para el formulario de contacto.

* librería vendorizada: `static/js/htmx.min.js`
* el formulario en `/contact` hace `hx-post="/contact"`
* el servidor valida y responde un fragmento HTML (`#contact-form-feedback`)
* sin HTMX, el mismo `POST` funciona (redirect o re-render)

El canal `mailto:info@dgsis.com` se mantiene en paralelo.

## Por qué HTMX aquí

Encaja con HTML-first: feedback parcial sin SPA ni framework frontend.

## Dónde no debe usarse

* navegación principal entre páginas
* sustituir SSR por un modelo SPA
* añadir complejidad “por si acaso”

## Principio

HTMX = mejora progresiva para necesidades reales. Hoy: contacto. Mañana: solo si aporta valor claro.

---

# 9. CSS Architecture

## Estrategia

Un único archivo CSS propio:

```
static/css/style.css
```

## Tokens actuales (`:root`)

* color: background, surface, text, muted, border
* tipografía: system font stack
* spacing: `--space-unit`
* layout: `--content-width`, `--section-max`

## Metodología

* selectores semánticos / BEM ligero (`.hero__action`, `.site-navigation__links`)
* composición por secciones
* media queries para responsive
* `prefers-reduced-motion` respetado

## Consistencia sin frameworks

La consistencia viene de:

* variables CSS compartidas
* patrones visuales repetidos (items con border-top)
* reglas del design system documental

No se usan Bootstrap, Tailwind ni librerías CSS externas.

---

# 10. JavaScript Policy

## Estado actual

No hay JavaScript en el repositorio.

## Cuándo estaría permitido

Solo si resuelve un problema concreto que HTML/CSS (y eventualmente HTMX) no puedan cubrir de forma razonable.

## Qué evitar

* frameworks SPA
* JS por defecto en cada página
* animaciones o tracking que no aporten a la experiencia de ingeniería

## Principio

JavaScript debe existir por necesidad, no por hábito.

---

# 11. Dependency Philosophy

## Regla

Preferir la biblioteca estándar. Cada dependencia nueva debe justificar:

1. el problema concreto
2. el beneficio medible
3. el coste de mantenimiento/actualización

## Estado

`go.mod` no declara dependencias de terceros.

## ¿Cuándo una dependencia está justificada?

Cuando la alternativa standard-library sea claramente peor en seguridad, corrección o mantenibilidad — y el alcance del problema sea real, no hipotético.

---

# 12. Testing Strategy

## Estado actual

No hay suites automatizadas en el repositorio todavía.

## Estrategia prevista (proporcional al tamaño)

| Tipo | Qué probar | Prioridad |
| --- | --- | --- |
| HTTP smoke | rutas principales → 200; desconocidas → 404; `/health` JSON | Alta |
| Template/render | handlers devuelven título/contenido esperado | Media |
| Unit | `getEnv` y helpers pequeños si crecen | Baja hoy |
| E2E browser | solo si la interacción crece (formularios/HTMX) | Futuro |

## Qué no necesita pruebas todavía

* lógica inexistente de dominio
* snapshots visuales pesados sin diseño final de tokens

La validación actual se hace por ejecución local, revisión HTML y mediciones Lighthouse/a11y cuando corresponde.

---

# 13. Security Considerations

## Prácticas actuales

* HTML escapado por defecto vía `html/template`
* sin evaluación de input de usuario (sin formularios server-side)
* contenedor con usuario no root (`appuser`)
* health endpoint sin datos sensibles
* configuración por entorno (`PORT`)
* headers HTTP: `X-Content-Type-Options`, `Referrer-Policy`, `X-Frame-Options`, `Permissions-Policy`
* timeouts en `http.Server` (`ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`)

## Límites actuales / preparación futura

* TLS y HSTS pertenecen al reverse proxy / edge (no implementado en este repo)
* si se añade formulario: validación, límites de tamaño, CSRF según diseño, y nunca loguear PII innecesaria
* secretos futuros nunca en el repo; solo entorno/secret manager

## Exposición

La aplicación sirve contenido público estático/SSR. No hay autenticación ni datos privados.

---

# 14. Performance Principles

## Objetivos

* HTML pequeño y legible
* CSS único y liviano
* cero JS en el camino crítico
* tipografías de sistema (sin webfonts)
* imágenes mínimas (favicon SVG)

## Filosofía Lighthouse

Priorizar puntuaciones altas por arquitectura simple, no por micro-optimizaciones opacas. Mejoras deben ser comprensibles.

## Assets

Servidos directamente por `FileServer` con `Cache-Control` largo para `/static/`.  
HTML y CSS pueden comprimirse con gzip cuando el cliente envía `Accept-Encoding: gzip`.  
No hay bundler ni minificador en el pipeline actual.

## Renderizado

SSR síncrono en proceso; sin hidratación cliente.

---

# 15. Deployment Architecture

## Implementado actualmente

### Local

```bash
go run ./cmd/web
```

### Docker build

```bash
docker build .
```

### Docker run

```bash
docker run -p 8080:8080 <image>
```

### Imagen

* stage build: `golang:1.26-alpine`
* stage runtime: `alpine:latest`
* binario `dgsis-website`
* copia `templates/` y `static/`
* `EXPOSE 8080`
* `USER appuser`

## Preparado para futuro (no implementado)

* CI/CD
* reverse proxy (TLS, compresión, headers)
* observabilidad avanzada (métricas/tracing)
* múltiples entornos con config externa ampliada

---

# 16. Evolution Strategy

## Cómo añadir nuevas páginas

1. Crear `templates/pages/<name>.html` con `{{define "content"}}`
2. Parsear con `mustParsePage`
3. Registrar ruta y `pageData`
4. Enlazar en navegación solo si la IA lo justifica

## Formularios

Cuando exista necesidad real:

* preferir HTML + handler Go
* evaluar HTMX para feedback parcial
* no introducir SPA

## Contenido dinámico

Solo si aparece una fuente de datos real. Hoy el contenido es estático en templates.

## Blog / i18n

Descartados o futuros según UX Strategy. No forman parte de la arquitectura activa.

## Regla de crecimiento

Ampliar el sistema sin romper la simplicidad inicial. Si una feature exige muchas capas nuevas, revisar primero si la feature es necesaria.

---

# 17. Architecture Constraints

## Decisiones que NO deben cambiarse sin razón fuerte

1. **No** introducir frameworks frontend SPA (React/Vue/Svelte apps) para este sitio.
2. **No** introducir frameworks CSS externos que oculten el sistema visual.
3. **No** crear abstracciones de dominio/services sin lógica real.
4. Mantener **Go standard library** mientras sea suficiente.
5. Mantener **HTML first** y SSR como modelo principal.
6. **No** convertir el sitio en plataforma autenticada sin un producto que lo justifique.
7. HTMX solo como mejora progresiva, nunca como base SPA.

---

# 18. Approved Decisions vs Future Possibilities

## Decisiones aprobadas e implementadas

* Go + `net/http` + `html/template`
* SSR multi-página: Home, Engineering, Work, Contact
* CSS propio con tokens
* Docker multi-stage
* logging `slog`
* health endpoint
* metadata básica + accesibilidad estructural
* contacto por `mailto:` **y** formulario HTMX con validación server-side
* logo/isotipo SVG hexagonal propio en header, footer y favicon
* case studies ilustrativos (marcados como placeholder)

## Decisiones aprobadas pero no implementadas aún

* entrega real de email (SMTP/provider) — hoy el form valida y registra en logs
* case studies de clientes verificados
* documento de status (`05`) actualizado tras cada incremento

## Posibles evoluciones futuras

* reverse proxy / TLS termination
* CI/CD
* formulario de contacto con validación
* división del CSS si crece materialmente
* tests HTTP automatizados
* headers de seguridad endurecidos de forma sistemática

---

# 19. Final Architecture Statement

En el ecosistema DGSIS, el software se construye para permanecer comprensible.

La arquitectura correcta es la que resuelve el problema real con la menor complejidad sostenible: decisiones explícitas, dependencias justificadas, HTML como base y un runtime que otro ingeniero pueda operar años después.

Esta web no existe para exhibir tecnología. Existe para demostrar criterio.

> Build systems that remain simple, understandable and maintainable.

---

# Relación con otros documentos

| Documento | Relación |
| --- | --- |
| `01_BRAND_BOOK.md` | Identidad y tono que la implementación debe reflejar |
| `02_UX_STRATEGY.md` | IA y páginas justificadas |
| `03_DESIGN_SYSTEM.md` | Principios visuales y CSS propio |
| `05_IMPLEMENTATION_STATUS.md` | Estado de avance y siguiente fase |

Si este documento y el código divergen, **el código desplegable manda** y este archivo debe actualizarse en el mismo cambio.

---

**Fin del documento**
