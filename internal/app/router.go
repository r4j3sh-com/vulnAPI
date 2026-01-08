package app

import "net/http"

func routes(store *Store, cfg Config) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", staticHandler()))
	mux.HandleFunc("/", indexPageHandler())
	mux.HandleFunc("/about", aboutPageHandler())
	mux.HandleFunc("/pricing", pricingPageHandler())
	mux.HandleFunc("/portal", portalPageHandler())
	mux.HandleFunc("/docs", docsPageHandler())
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/auth/login", loginHandler(store))
	mux.HandleFunc("/api/auth/reset", passwordResetHandler(store))
	mux.HandleFunc("/api/auth/dashboard", dashboardHandler(store))
	mux.HandleFunc("/api/users", usersHandler(store))
	mux.HandleFunc("/api/users/", userHandler(store))
	mux.HandleFunc("/api/orders", ordersHandler(store))
	mux.HandleFunc("/api/orders/", orderHandler(store))
	mux.HandleFunc("/api/catalog", productsHandler(store))
	mux.HandleFunc("/api/catalog/search", searchHandler(store))
	mux.HandleFunc("/api/reports", reportHandler())
	mux.HandleFunc("/api/ops/metrics", adminStatsHandler(store))
	mux.HandleFunc("/api/ops/roles/promote/", promoteHandler(store))
	mux.HandleFunc("/api/ops/ping", pingHandler())
	mux.HandleFunc("/api/ops/files", fileReadHandler())

	mux.HandleFunc("/api/v1/users", v1UsersHandler(store))
	mux.HandleFunc("/api/v2/users", v2UsersHandler(store))

	mux.HandleFunc("/ops/config", debugConfigHandler(cfg))

	return mux
}
