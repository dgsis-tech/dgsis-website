package main

import (
	"net/http"
	"strings"
)

const (
	langCookieName = "dgsis_lang"
	langES         = "es"
	langEN         = "en"
	langCookieMax  = 365 * 24 * 60 * 60
)

// Msg holds UI copy for one language. Templates read fields via {{.Msg.*}}.
type Msg struct {
	SkipToContent  string
	NavAria        string
	NavHome        string
	NavEngineering string
	NavWork        string
	NavContact     string

	LangLabel   string
	ThemeLabel  string
	ThemeLight  string
	ThemeDark   string
	ThemeSystem string
	PrefsLabel  string

	FooterTagline   string
	FooterCopyright string

	HomeTitle        string
	HomeDescription  string
	HomeEyebrow      string
	HomeH1           string
	HomeLead         string
	HomeCTA          string
	HomeCTASecondary string
	HomeApproachH2   string
	HomeApproachLead string
	HomeSimplicityH3 string
	HomeSimplicityP  string
	HomeMaintainH3   string
	HomeMaintainP    string
	HomeJudgmentH3   string
	HomeJudgmentP    string
	HomeEngH2        string
	HomeEngLead      string
	HomeBackendH3    string
	HomeBackendP     string
	HomeArchH3       string
	HomeArchP        string
	HomePlatformsH3  string
	HomePlatformsP   string
	HomeEngLink      string
	HomeWorkH2       string
	HomeWorkLead     string
	HomeProvisional  string
	HomeWork1H3      string
	HomeWork1P       string
	HomeWork2H3      string
	HomeWork2P       string
	HomeWork3H3      string
	HomeWork3P       string
	HomeWorkLink     string
	HomeContactH2    string
	HomeContactLead  string

	EngTitle       string
	EngDescription string
	EngEyebrow     string
	EngH1          string
	EngLead        string
	EngFocusH2     string
	EngFocusLead   string
	EngBackendH3   string
	EngBackendP    string
	EngArchH3      string
	EngArchP       string
	EngPlatformsH3 string
	EngPlatformsP  string
	EngIntegrH3    string
	EngIntegrP     string
	EngAutoH3      string
	EngAutoP       string
	EngSaaSH3      string
	EngSaaSP       string
	EngMethodH2    string
	EngMethodLead  string
	EngStep1H3     string
	EngStep1P      string
	EngStep2H3     string
	EngStep2P      string
	EngStep3H3     string
	EngStep3P      string
	EngContactH2   string
	EngContactLead string
	EngCTA         string

	WorkTitle       string
	WorkDescription string
	WorkEyebrow     string
	WorkH1          string
	WorkLead        string
	WorkProvisional string
	WorkCaseLabel   string
	WorkProblem     string
	WorkDecisions   string
	WorkSolution    string
	WorkResult      string
	WorkCase1H2     string
	WorkCase1Sum    string
	WorkCase1Prob   string
	WorkCase1Dec    string
	WorkCase1Sol    string
	WorkCase1Res    string
	WorkCase2H2     string
	WorkCase2Sum    string
	WorkCase2Prob   string
	WorkCase2Dec    string
	WorkCase2Sol    string
	WorkCase2Res    string
	WorkCase3H2     string
	WorkCase3Sum    string
	WorkCase3Prob   string
	WorkCase3Dec    string
	WorkCase3Sol    string
	WorkCase3Res    string
	WorkContactH2   string
	WorkContactLead string
	WorkCTA         string

	ContactTitle       string
	ContactDescription string
	ContactEyebrow     string
	ContactH1          string
	ContactLead        string
	ContactConvH2      string
	ContactConvLead    string
	ContactMailto      string
	ContactAvailLabel  string
	ContactAvailValue  string
	ContactModeLabel   string
	ContactModeValue   string
	ContactStackLabel  string
	ContactStackValue  string
	ContactName        string
	ContactEmail       string
	ContactCompany     string
	ContactOptional    string
	ContactMessage     string
	ContactPlaceholder string
	ContactSubmit      string
	ContactNote        string
	ContactOK          string
	ContactOKShort     string
	ContactErrRead     string
	ContactErrName     string
	ContactErrEmail    string
	ContactErrCompany  string
	ContactErrMessage  string
}

func messages(lang string) Msg {
	if lang == langEN {
		return msgEN
	}
	return msgES
}

func normalizeLang(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return ""
	}
	primary := v
	if i := strings.IndexByte(v, '-'); i >= 0 {
		primary = v[:i]
	} else if i := strings.IndexByte(v, '_'); i >= 0 {
		primary = v[:i]
	}
	switch primary {
	case langES:
		return langES
	case langEN:
		return langEN
	default:
		return ""
	}
}

// resolveLang picks language: ?lang= → cookie → Accept-Language → ES.
// When ?lang= is present and valid, the cookie is set and the query is stripped via redirect.
func resolveLang(w http.ResponseWriter, r *http.Request) (lang string, redirected bool) {
	if q := normalizeLang(r.URL.Query().Get("lang")); q != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     langCookieName,
			Value:    q,
			Path:     "/",
			MaxAge:   langCookieMax,
			SameSite: http.SameSiteLaxMode,
			HttpOnly: false,
		})

		values := r.URL.Query()
		values.Del("lang")
		r.URL.RawQuery = values.Encode()
		target := r.URL.Path
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusSeeOther)
		return q, true
	}

	if c, err := r.Cookie(langCookieName); err == nil {
		if v := normalizeLang(c.Value); v != "" {
			return v, false
		}
	}

	for _, part := range strings.Split(r.Header.Get("Accept-Language"), ",") {
		token := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		if v := normalizeLang(token); v != "" {
			return v, false
		}
	}

	// Default: Spanish (company primary audience / brand docs). Browser EN still wins above.
	return langES, false
}

var msgEN = Msg{
	SkipToContent:  "Skip to content",
	NavAria:        "Main navigation",
	NavHome:        "Home",
	NavEngineering: "Engineering",
	NavWork:        "Work",
	NavContact:     "Contact",

	LangLabel:   "Language",
	ThemeLabel:  "Theme",
	ThemeLight:  "Light",
	ThemeDark:   "Dark",
	ThemeSystem: "System",
	PrefsLabel:  "Language and theme",

	FooterTagline:   "Engineering software that lasts.",
	FooterCopyright: "© DGSIS LLC",

	HomeTitle:        "DGSIS — Engineering software that lasts.",
	HomeDescription:  "DGSIS is a software engineering company that builds reliable digital systems designed to evolve, scale and remain maintainable for years.",
	HomeEyebrow:      "Software engineering company",
	HomeH1:           "Engineering software that lasts.",
	HomeLead:         "We build reliable digital systems designed to evolve, scale and remain maintainable for years.",
	HomeCTA:          "Start a conversation",
	HomeCTASecondary: "Explore engineering",
	HomeApproachH2:   "Engineering with intention.",
	HomeApproachLead: "Good software is not defined by how much technology it uses, but by the quality of the decisions behind it.",
	HomeSimplicityH3: "Simplicity",
	HomeSimplicityP:  "The simplest solution that solves the problem correctly is usually the best one.",
	HomeMaintainH3:   "Maintainability",
	HomeMaintainP:    "Software should remain understandable and adaptable years after it is built.",
	HomeJudgmentH3:   "Technical judgment",
	HomeJudgmentP:    "Every technical decision should have a clear reason behind it.",
	HomeEngH2:        "Engineering systems that scale with purpose.",
	HomeEngLead:      "We design and build digital systems focused on reliability, performance and long-term maintainability.",
	HomeBackendH3:    "Backend Systems",
	HomeBackendP:     "Reliable services and APIs built with simplicity, clarity and operational stability in mind.",
	HomeArchH3:       "Software Architecture",
	HomeArchP:        "Systems designed around clear boundaries, maintainable structures and thoughtful decisions.",
	HomePlatformsH3:  "Digital Platforms",
	HomePlatformsP:   "Complete platforms built to support real users, business processes and future growth.",
	HomeEngLink:      "Read more about our engineering approach",
	HomeWorkH2:       "Selected engineering work.",
	HomeWorkLead:     "Every system we build starts with understanding the problem, making deliberate technical decisions and creating foundations that can evolve over time.",
	HomeProvisional:  "Illustrative placeholders — not verified client engagements.",
	HomeWork1H3:      "Operational ledger",
	HomeWork1P:       "Clear domain boundaries and operable APIs for growing product operations.",
	HomeWork2H3:      "Learning platform foundation",
	HomeWork2P:       "Architecture shaped for content growth without structural debt.",
	HomeWork3H3:      "Process automation",
	HomeWork3P:       "Restrained automation with visible failure modes and explicit ownership.",
	HomeWorkLink:     "View case studies",
	HomeContactH2:    "Let's build something that lasts.",
	HomeContactLead:  "Have a project that requires thoughtful engineering? Let's discuss how we can build a solid foundation together.",

	EngTitle:       "Engineering — DGSIS",
	EngDescription: "DGSIS designs and builds backend systems, software architecture and digital platforms focused on reliability and long-term maintainability.",
	EngEyebrow:     "Engineering",
	EngH1:          "Systems built with clear technical judgment.",
	EngLead:        "DGSIS designs and builds digital systems where technology serves the problem. We choose tools for clarity, reliability and the ability to evolve — not for novelty.",
	EngFocusH2:     "Where we focus.",
	EngFocusLead:   "Our work concentrates on the foundations that make platforms stable today and adaptable tomorrow.",
	EngBackendH3:   "Backend systems",
	EngBackendP:    "Services and APIs designed for operational stability, clear contracts and long-term ownership.",
	EngArchH3:      "Software architecture",
	EngArchP:       "Boundaries, responsibilities and structures chosen for the real problem — not the most complex diagram.",
	EngPlatformsH3: "Digital platforms",
	EngPlatformsP:  "End-to-end platforms that support users, business processes and controlled growth over time.",
	EngIntegrH3:    "Integrations",
	EngIntegrP:     "Connections between systems designed to reduce fragility and keep operational intent visible.",
	EngAutoH3:      "Process automation",
	EngAutoP:       "Software that replaces fragile manual workflows with systems that remain understandable as they grow.",
	EngSaaSH3:      "SaaS foundations",
	EngSaaSP:       "Product platforms built to stay maintainable as usage, teams and requirements expand.",
	EngMethodH2:    "How we build.",
	EngMethodLead:  "We do not sell technology lists. We apply engineering judgment so systems remain simple where possible and deliberate where complexity is required.",
	EngStep1H3:     "Understand the problem",
	EngStep1P:      "Before choosing tools, we clarify constraints, users and the decisions that will matter years later.",
	EngStep2H3:     "Prefer lasting structures",
	EngStep2P:      "We favor maintainable designs over temporary shortcuts that become expensive to reverse.",
	EngStep3H3:     "Evolve with control",
	EngStep3P:      "Systems should change without losing stability. Growth is planned as part of the architecture, not an afterthought.",
	EngContactH2:   "Need a durable technical foundation?",
	EngContactLead: "If you are building a product, modernizing a process or strengthening an existing platform, we can help you decide and build with intention.",
	EngCTA:         "Start a conversation",

	WorkTitle:       "Work — DGSIS",
	WorkDescription: "Selected engineering work from DGSIS. Case studies focused on problems, decisions, solutions and lasting results.",
	WorkEyebrow:     "Work",
	WorkH1:          "Selected engineering work.",
	WorkLead:        "Each case follows the same structure: problem, decisions, solution and result. These examples illustrate how DGSIS thinks about durable systems.",
	WorkProvisional: "Illustrative placeholders — not verified client engagements.",
	WorkCaseLabel:   "Placeholder case study",
	WorkProblem:     "Problem",
	WorkDecisions:   "Decisions",
	WorkSolution:    "Solution",
	WorkResult:      "Result",
	WorkCase1H2:     "Operational ledger for a growing product company",
	WorkCase1Sum:    "A product team needed a reliable core for billing-related operations without locking the business into an opaque platform.",
	WorkCase1Prob:   "Manual reconciliation and loosely coupled scripts made financial operations fragile as volume grew. Changes were slow and risky because knowledge lived in a few people.",
	WorkCase1Dec:    "Model clear domain boundaries, keep the service surface small, and prefer explicit workflows over hidden side effects. Choose boring, operable building blocks over a complex event mesh.",
	WorkCase1Sol:    "A focused backend service with stable APIs, audited state transitions and operational visibility. Integrations were isolated so the core ledger could evolve independently.",
	WorkCase1Res:    "Operations became repeatable, onboarding new engineers got simpler, and the company could change pricing rules without rewriting surrounding systems.",
	WorkCase2H2:     "Learning platform foundation for continuous growth",
	WorkCase2Sum:    "An education business needed a digital platform that could support more users and content without accumulating structural debt.",
	WorkCase2Prob:   "The existing application mixed presentation, business rules and integrations. Every feature release increased the cost of the next change.",
	WorkCase2Dec:    "Separate the platform into maintainable layers, keep the HTML-facing experience simple, and design data models for long-term content growth rather than short-term screens.",
	WorkCase2Sol:    "A clearer architecture for users, courses and access rules, with backend services shaped around real workflows instead of UI convenience.",
	WorkCase2Res:    "The team could ship platform improvements with less regression risk, and the system remained understandable as the catalog and audience expanded.",
	WorkCase3H2:     "Process automation for internal operations",
	WorkCase3Sum:    "A company digitalizing internal processes needed software that reduced manual work without creating a fragile integration tangle.",
	WorkCase3Prob:   "Critical workflows depended on spreadsheets and ad-hoc tools. Errors were hard to detect, and every new system increased operational uncertainty.",
	WorkCase3Dec:    "Automate only the steps that create durable value, keep human checkpoints where judgment matters, and make every integration ownership-explicit.",
	WorkCase3Sol:    "A restrained automation layer with clear inputs/outputs, recoverable jobs and logs that operators can trust during incidents.",
	WorkCase3Res:    "Cycle time dropped on routine work, failure modes became visible, and the organization could extend automation without rewriting the foundation.",
	WorkContactH2:   "Building something that needs to last?",
	WorkContactLead: "Tell us about the system you need to create or strengthen. We can discuss the technical path with the same care we apply to delivery.",
	WorkCTA:         "Start a conversation",

	ContactTitle:       "Contact — DGSIS",
	ContactDescription: "Start a conversation with DGSIS about building digital systems with solid engineering foundations.",
	ContactEyebrow:     "Contact",
	ContactH1:          "Let's build something that lasts.",
	ContactLead:        "If you have a project that needs thoughtful engineering, we welcome a direct conversation about the problem, constraints and the foundation required.",
	ContactConvH2:      "Start a technical conversation.",
	ContactConvLead:    "Share context about your product, platform or process. Use the form below, or email us directly.",
	ContactMailto:      "Email info@dgsis.com",
	ContactAvailLabel:  "AVAILABILITY",
	ContactAvailValue:  "Open to new projects",
	ContactModeLabel:   "ENGAGEMENT",
	ContactModeValue:   "Remote, projects and retainer",
	ContactStackLabel:  "CORE STACK",
	ContactStackValue:  "Go · gRPC · PostgreSQL · HTMX",
	ContactName:        "Name",
	ContactEmail:       "Email",
	ContactCompany:     "Company",
	ContactOptional:    "(optional)",
	ContactMessage:     "Project context",
	ContactPlaceholder: "What are you building, what constraints matter, and what does success look like?",
	ContactSubmit:      "Send message",
	ContactNote:        "Prefer email? Write to info@dgsis.com. Form submissions are validated on the server and recorded for follow-up.",
	ContactOK:          "Thank you. Your message was received. You can also reach us at info@dgsis.com.",
	ContactOKShort:     "Thank you. Your message was received.",
	ContactErrRead:     "The message could not be read. Please try again or email info@dgsis.com.",
	ContactErrName:     "Please enter a name between 2 and 80 characters.",
	ContactErrEmail:    "Please enter a valid email address.",
	ContactErrCompany:  "Company name is too long.",
	ContactErrMessage:  "Please describe the project in 20–4000 characters.",
}

var msgES = Msg{
	SkipToContent:  "Saltar al contenido",
	NavAria:        "Navegación principal",
	NavHome:        "Inicio",
	NavEngineering: "Ingeniería",
	NavWork:        "Trabajo",
	NavContact:     "Contacto",

	LangLabel:   "Idioma",
	ThemeLabel:  "Tema",
	ThemeLight:  "Claro",
	ThemeDark:   "Oscuro",
	ThemeSystem: "Sistema",
	PrefsLabel:  "Idioma y tema",

	FooterTagline:   "Ingeniería de software que perdura.",
	FooterCopyright: "© DGSIS LLC",

	HomeTitle:        "DGSIS — Ingeniería de software que perdura.",
	HomeDescription:  "DGSIS es una empresa de ingeniería de software que construye sistemas digitales fiables, pensados para evolucionar, escalar y permanecer mantenibles durante años.",
	HomeEyebrow:      "Empresa de ingeniería de software",
	HomeH1:           "Ingeniería de software que perdura.",
	HomeLead:         "Construimos sistemas digitales fiables, diseñados para evolucionar, escalar y permanecer mantenibles durante años.",
	HomeCTA:          "Iniciar una conversación",
	HomeCTASecondary: "Explorar ingeniería",
	HomeApproachH2:   "Ingeniería con intención.",
	HomeApproachLead: "El buen software no se define por la cantidad de tecnología que usa, sino por la calidad de las decisiones que lo sostienen.",
	HomeSimplicityH3: "Simplicidad",
	HomeSimplicityP:  "La solución más simple que resuelve el problema correctamente suele ser la mejor.",
	HomeMaintainH3:   "Mantenibilidad",
	HomeMaintainP:    "El software debe seguir siendo comprensible y adaptable años después de construirse.",
	HomeJudgmentH3:   "Criterio técnico",
	HomeJudgmentP:    "Cada decisión técnica debe tener una razón clara.",
	HomeEngH2:        "Sistemas de ingeniería que escalan con propósito.",
	HomeEngLead:      "Diseñamos y construimos sistemas digitales centrados en fiabilidad, rendimiento y mantenibilidad a largo plazo.",
	HomeBackendH3:    "Sistemas backend",
	HomeBackendP:     "Servicios y APIs fiables, construidos con simplicidad, claridad y estabilidad operativa.",
	HomeArchH3:       "Arquitectura de software",
	HomeArchP:        "Sistemas diseñados con límites claros, estructuras mantenibles y decisiones fundamentadas.",
	HomePlatformsH3:  "Plataformas digitales",
	HomePlatformsP:   "Plataformas completas para usuarios reales, procesos de negocio y crecimiento futuro.",
	HomeEngLink:      "Leer más sobre nuestro enfoque de ingeniería",
	HomeWorkH2:       "Trabajo de ingeniería seleccionado.",
	HomeWorkLead:     "Cada sistema empieza por entender el problema, tomar decisiones técnicas deliberadas y crear bases que puedan evolucionar.",
	HomeProvisional:  "Ejemplos ilustrativos — no son encargos de cliente verificados.",
	HomeWork1H3:      "Ledger operativo",
	HomeWork1P:       "Límites de dominio claros y APIs operables para operaciones de producto en crecimiento.",
	HomeWork2H3:      "Base de plataforma de aprendizaje",
	HomeWork2P:       "Arquitectura preparada para crecer en contenido sin deuda estructural.",
	HomeWork3H3:      "Automatización de procesos",
	HomeWork3P:       "Automatización contenida, con fallos visibles y responsabilidad explícita.",
	HomeWorkLink:     "Ver casos de estudio",
	HomeContactH2:    "Construyamos algo que perdure.",
	HomeContactLead:  "¿Tienes un proyecto que requiere ingeniería con criterio? Hablemos de cómo construir una base sólida.",

	EngTitle:       "Ingeniería — DGSIS",
	EngDescription: "DGSIS diseña y construye sistemas backend, arquitectura de software y plataformas digitales centradas en fiabilidad y mantenibilidad a largo plazo.",
	EngEyebrow:     "Ingeniería",
	EngH1:          "Sistemas construidos con criterio técnico claro.",
	EngLead:        "DGSIS diseña y construye sistemas digitales donde la tecnología sirve al problema. Elegimos herramientas por claridad, fiabilidad y capacidad de evolucionar — no por novedad.",
	EngFocusH2:     "En qué nos centramos.",
	EngFocusLead:   "Nuestro trabajo se concentra en las bases que hacen las plataformas estables hoy y adaptables mañana.",
	EngBackendH3:   "Sistemas backend",
	EngBackendP:    "Servicios y APIs diseñados para estabilidad operativa, contratos claros y ownership a largo plazo.",
	EngArchH3:      "Arquitectura de software",
	EngArchP:       "Límites, responsabilidades y estructuras elegidos para el problema real — no para el diagrama más complejo.",
	EngPlatformsH3: "Plataformas digitales",
	EngPlatformsP:  "Plataformas de extremo a extremo que sostienen usuarios, procesos de negocio y crecimiento controlado.",
	EngIntegrH3:    "Integraciones",
	EngIntegrP:     "Conexiones entre sistemas pensadas para reducir fragilidad y mantener visible la intención operativa.",
	EngAutoH3:      "Automatización de procesos",
	EngAutoP:       "Software que sustituye flujos manuales frágiles por sistemas que siguen siendo comprensibles al crecer.",
	EngSaaSH3:      "Fundamentos SaaS",
	EngSaaSP:       "Plataformas de producto construidas para permanecer mantenibles cuando crecen uso, equipos y requisitos.",
	EngMethodH2:    "Cómo construimos.",
	EngMethodLead:  "No vendemos listas de tecnología. Aplicamos criterio de ingeniería para que los sistemas sean simples cuando es posible y deliberados cuando la complejidad es necesaria.",
	EngStep1H3:     "Entender el problema",
	EngStep1P:      "Antes de elegir herramientas, aclaramos restricciones, usuarios y las decisiones que importarán años después.",
	EngStep2H3:     "Preferir estructuras duraderas",
	EngStep2P:      "Priorizamos diseños mantenibles frente a atajos temporales caros de revertir.",
	EngStep3H3:     "Evolucionar con control",
	EngStep3P:      "Los sistemas deben cambiar sin perder estabilidad. El crecimiento se planifica como parte de la arquitectura, no como un añadido.",
	EngContactH2:   "¿Necesitas una base técnica duradera?",
	EngContactLead: "Si estás construyendo un producto, modernizando un proceso o reforzando una plataforma existente, podemos ayudarte a decidir y construir con intención.",
	EngCTA:         "Iniciar una conversación",

	WorkTitle:       "Trabajo — DGSIS",
	WorkDescription: "Trabajo de ingeniería seleccionado de DGSIS. Casos centrados en problema, decisiones, solución y resultado duradero.",
	WorkEyebrow:     "Trabajo",
	WorkH1:          "Trabajo de ingeniería seleccionado.",
	WorkLead:        "Cada caso sigue la misma estructura: problema, decisiones, solución y resultado. Estos ejemplos ilustran cómo DGSIS piensa sistemas duraderos.",
	WorkProvisional: "Ejemplos ilustrativos — no son encargos de cliente verificados.",
	WorkCaseLabel:   "Caso de estudio provisional",
	WorkProblem:     "Problema",
	WorkDecisions:   "Decisiones",
	WorkSolution:    "Solución",
	WorkResult:      "Resultado",
	WorkCase1H2:     "Ledger operativo para una empresa de producto en crecimiento",
	WorkCase1Sum:    "Un equipo de producto necesitaba un núcleo fiable para operaciones de facturación sin quedar atrapado en una plataforma opaca.",
	WorkCase1Prob:   "La reconciliación manual y scripts débilmente acoplados hacían frágiles las operaciones financieras al crecer el volumen. Los cambios eran lentos y arriesgados porque el conocimiento vivía en pocas personas.",
	WorkCase1Dec:    "Modelar límites de dominio claros, mantener pequeña la superficie del servicio y preferir flujos explícitos frente a efectos ocultos. Elegir bloques operables y sobrios frente a una malla de eventos compleja.",
	WorkCase1Sol:    "Un servicio backend enfocado con APIs estables, transiciones de estado auditadas y visibilidad operativa. Las integraciones se aislaron para que el ledger pudiera evolucionar de forma independiente.",
	WorkCase1Res:    "Las operaciones se volvieron repetibles, incorporar ingenieros fue más simple y la empresa pudo cambiar reglas de precios sin reescribir los sistemas alrededor.",
	WorkCase2H2:     "Base de plataforma de aprendizaje para crecimiento continuo",
	WorkCase2Sum:    "Un negocio educativo necesitaba una plataforma digital capaz de soportar más usuarios y contenido sin acumular deuda estructural.",
	WorkCase2Prob:   "La aplicación existente mezclaba presentación, reglas de negocio e integraciones. Cada release encarecía el siguiente cambio.",
	WorkCase2Dec:    "Separar la plataforma en capas mantenibles, mantener simple la experiencia HTML y diseñar modelos de datos para crecimiento de contenido a largo plazo, no para pantallas a corto plazo.",
	WorkCase2Sol:    "Una arquitectura más clara para usuarios, cursos y reglas de acceso, con servicios backend orientados a flujos reales en lugar de conveniencia de UI.",
	WorkCase2Res:    "El equipo pudo entregar mejoras de plataforma con menos riesgo de regresión, y el sistema siguió siendo comprensible al crecer catálogo y audiencia.",
	WorkCase3H2:     "Automatización de procesos para operaciones internas",
	WorkCase3Sum:    "Una empresa que digitalizaba procesos internos necesitaba software que redujera trabajo manual sin crear un entramado frágil de integraciones.",
	WorkCase3Prob:   "Flujos críticos dependían de hojas de cálculo y herramientas ad hoc. Los errores eran difíciles de detectar y cada sistema nuevo aumentaba la incertidumbre operativa.",
	WorkCase3Dec:    "Automatizar solo los pasos que aportan valor duradero, mantener puntos de control humanos donde importa el criterio y hacer explícita la responsabilidad de cada integración.",
	WorkCase3Sol:    "Una capa de automatización contenida, con entradas/salidas claras, trabajos recuperables y logs en los que los operadores pueden confiar durante incidentes.",
	WorkCase3Res:    "Bajó el tiempo de ciclo en trabajo rutinario, los modos de fallo se hicieron visibles y la organización pudo extender la automatización sin reescribir la base.",
	WorkContactH2:   "¿Estás construyendo algo que debe perdurar?",
	WorkContactLead: "Cuéntanos el sistema que necesitas crear o reforzar. Podemos hablar del camino técnico con el mismo cuidado que aplicamos a la entrega.",
	WorkCTA:         "Iniciar una conversación",

	ContactTitle:       "Contacto — DGSIS",
	ContactDescription: "Inicia una conversación con DGSIS sobre construir sistemas digitales con bases sólidas de ingeniería.",
	ContactEyebrow:     "Contacto",
	ContactH1:          "Construyamos algo que perdure.",
	ContactLead:        "Si tienes un proyecto que necesita ingeniería con criterio, te invitamos a una conversación directa sobre el problema, las restricciones y la base requerida.",
	ContactConvH2:      "Inicia una conversación técnica.",
	ContactConvLead:    "Comparte contexto sobre tu producto, plataforma o proceso. Usa el formulario o escríbenos por email.",
	ContactMailto:      "Email info@dgsis.com",
	ContactAvailLabel:  "DISPONIBILIDAD",
	ContactAvailValue:  "Abiertos a nuevos proyectos",
	ContactModeLabel:   "MODALIDAD",
	ContactModeValue:   "Remoto, proyectos y retainer",
	ContactStackLabel:  "STACK PRINCIPAL",
	ContactStackValue:  "Go · gRPC · PostgreSQL · HTMX",
	ContactName:        "Nombre",
	ContactEmail:       "Email",
	ContactCompany:     "Empresa",
	ContactOptional:    "(opcional)",
	ContactMessage:     "Contexto del proyecto",
	ContactPlaceholder: "¿Qué estás construyendo, qué restricciones importan y cómo se ve el éxito?",
	ContactSubmit:      "Enviar mensaje",
	ContactNote:        "¿Prefieres email? Escribe a info@dgsis.com. Los envíos del formulario se validan en el servidor y se registran para seguimiento.",
	ContactOK:          "Gracias. Hemos recibido tu mensaje. También puedes escribirnos a info@dgsis.com.",
	ContactOKShort:     "Gracias. Hemos recibido tu mensaje.",
	ContactErrRead:     "No se pudo leer el mensaje. Inténtalo de nuevo o escribe a info@dgsis.com.",
	ContactErrName:     "Introduce un nombre de entre 2 y 80 caracteres.",
	ContactErrEmail:    "Introduce una dirección de email válida.",
	ContactErrCompany:  "El nombre de la empresa es demasiado largo.",
	ContactErrMessage:  "Describe el proyecto en 20–4000 caracteres.",
}
