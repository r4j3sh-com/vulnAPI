package app

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

//go:embed web/templates/*.html web/static/**
var siteFS embed.FS

var siteTemplates = template.Must(template.New("site").Funcs(template.FuncMap{
	"year": func() int { return time.Now().Year() },
}).ParseFS(siteFS, "web/templates/*.html"))

type PageData struct {
	Title         string
	Description   string
	ContentHTML   template.HTML
	HeroTag       string
	HeroTitle     string
	HeroLead      string
	PrimaryCTA    string
	PrimaryLink   string
	SecondaryCTA  string
	SecondaryLink string
	Features      []Feature
	Metrics       []Metric
	Plans         []Plan
	Testimonials  []Testimonial
	Docs          []DocItem
}

type Feature struct {
	Title string
	Body  string
}

type Metric struct {
	Label string
	Value string
}

type Plan struct {
	Name        string
	Price       string
	Tagline     string
	Highlights  []string
	AccentClass string
}

type Testimonial struct {
	Quote string
	Name  string
	Role  string
}

type DocItem struct {
	Title   string
	Body    string
	Example string
}

func renderPage(w http.ResponseWriter, data PageData, contentTemplate string) {
	var buf bytes.Buffer
	if err := siteTemplates.ExecuteTemplate(&buf, contentTemplate, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	data.ContentHTML = template.HTML(buf.String())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := siteTemplates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func staticHandler() http.Handler {
	staticFS, err := fs.Sub(siteFS, "web/static")
	if err != nil {
		return http.NotFoundHandler()
	}
	return http.FileServer(http.FS(staticFS))
}

func indexPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:         "Harbor Cloud - Retail telemetry built for operators",
			Description:   "A modern telemetry platform for retail operations, order flow, and device health.",
			HeroTag:       "Retail telemetry for modern operations",
			HeroTitle:     "Operate every storefront like a control room.",
			HeroLead:      "Harbor Cloud aggregates orders, devices, and customer signals into one API and live control panel. Built for teams who need answers in seconds.",
			PrimaryCTA:    "Launch the demo",
			PrimaryLink:   "/portal",
			SecondaryCTA:  "Read the API guide",
			SecondaryLink: "/docs",
			Features: []Feature{
				{Title: "Unified order graph", Body: "Correlate customers, orders, and inventory across regions with a single API surface."},
				{Title: "Edge visibility", Body: "Track device health, latency, and outages with diagnostics routed to your operations team."},
				{Title: "Partner access", Body: "Issue keys to third-party vendors with scoped access to critical workflows."},
				{Title: "Fast onboarding", Body: "Bring a new store online in under a day with pre-built connectors."},
			},
			Metrics: []Metric{
				{Label: "Median response", Value: "42ms"},
				{Label: "Connected stores", Value: "128"},
				{Label: "Daily requests", Value: "2.3M"},
			},
			Testimonials: []Testimonial{
				{Quote: "Harbor turned our incident response into a live dashboard our whole team understands.", Name: "Nora Patel", Role: "Ops Lead, Coastline Retail"},
				{Quote: "Our vendors finally have an API they can integrate without three weeks of calls.", Name: "Marco Ruiz", Role: "Platform Director, Dune Markets"},
			},
		}
		renderPage(w, data, "index-content")
	}
}

func aboutPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:         "About Harbor Cloud",
			Description:   "Built by operators for operators, Harbor Cloud keeps modern retail resilient.",
			HeroTag:       "A platform shaped by the store floor",
			HeroTitle:     "We obsess over uptime, not dashboards.",
			HeroLead:      "Harbor started inside a regional retail chain. We built the tooling we needed to keep stores online, and turned it into a platform teams can trust.",
			PrimaryCTA:    "Meet the team",
			PrimaryLink:   "/portal",
			SecondaryCTA:  "View pricing",
			SecondaryLink: "/pricing",
			Features: []Feature{
				{Title: "Incident-ready", Body: "Metrics flow straight into alerting so you see what customers feel, not just what servers log."},
				{Title: "Operator-first", Body: "Every screen is built for the daily rhythm of store managers and on-call engineers."},
				{Title: "Privacy by design", Body: "Regional data routing and granular access controls keep sensitive info in the right hands."},
			},
		}
		renderPage(w, data, "about-content")
	}
}

func pricingPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:         "Harbor Cloud Pricing",
			Description:   "Flexible plans for growing retail operations.",
			HeroTag:       "Simple pricing",
			HeroTitle:     "Scale visibility without scaling overhead.",
			HeroLead:      "Pick a plan that fits the number of storefronts, integrations, and partners you manage today.",
			PrimaryCTA:    "Request access",
			PrimaryLink:   "/portal",
			SecondaryCTA:  "Compare features",
			SecondaryLink: "/docs",
			Plans: []Plan{
				{Name: "Starter", Price: "$249", Tagline: "Per month, 10 locations", AccentClass: "plan-core", Highlights: []string{"Core telemetry", "Basic alerts", "Partner API"}},
				{Name: "Operations", Price: "$799", Tagline: "Per month, 50 locations", AccentClass: "plan-accent", Highlights: []string{"Advanced diagnostics", "Incident workflows", "Dedicated support"}},
				{Name: "Enterprise", Price: "Custom", Tagline: "Regional & global", AccentClass: "plan-elite", Highlights: []string{"Hybrid routing", "Audit exports", "Custom SLAs"}},
			},
		}
		renderPage(w, data, "pricing-content")
	}
}

func portalPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:         "Developer Portal",
			Description:   "Sign in to manage keys, apps, and webhook deliveries.",
			HeroTag:       "Developer portal",
			HeroTitle:     "Ship integrations faster with Harbor Cloud.",
			HeroLead:      "Access keys, test webhooks, and explore operational data using the Harbor Cloud API suite.",
			PrimaryCTA:    "View API docs",
			PrimaryLink:   "/docs",
			SecondaryCTA:  "Contact support",
			SecondaryLink: "/about",
			Features: []Feature{
				{Title: "Issue API keys", Body: "Generate keys for internal apps, vendors, and third-party integrators."},
				{Title: "Sandbox environment", Body: "Validate integrations before pushing to production."},
				{Title: "Webhook previews", Body: "Inspect payloads from shipping, billing, and inventory pipelines."},
			},
		}
		renderPage(w, data, "portal-content")
	}
}

func docsPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:         "Harbor Cloud API",
			Description:   "Reference endpoints for order flow, customer profiles, and operations diagnostics.",
			HeroTag:       "API reference",
			HeroTitle:     "Operational data where you need it.",
			HeroLead:      "Use the Harbor Cloud API to query orders, manage profiles, and monitor infrastructure in real time.",
			PrimaryCTA:    "Get an API key",
			PrimaryLink:   "/portal",
			SecondaryCTA:  "See examples",
			SecondaryLink: "#examples",
			Docs: []DocItem{
				{Title: "Authenticate", Body: "Exchange credentials for a session token before calling protected endpoints.", Example: "POST /api/auth/login"},
				{Title: "Orders", Body: "Retrieve recent orders or create new shipments for fulfillment.", Example: "GET /api/orders"},
				{Title: "Catalog", Body: "Search product inventory, pricing, and vendor metadata.", Example: "GET /api/catalog/search?q=flow"},
				{Title: "Operations", Body: "Access diagnostics and metrics for store infrastructure.", Example: "GET /api/ops/metrics"},
			},
		}
		renderPage(w, data, "docs-content")
	}
}
