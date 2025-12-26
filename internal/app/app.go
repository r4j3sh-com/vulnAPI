package app

import (
	"log"
	"net/http"
)

func Run() {
	cfg := LoadConfig()
	db, err := InitDatabase(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	store := NewStore(db)

	handler := withCORS(routes(store, cfg))

	addr := ":" + cfg.Port
	log.Printf("harbor-api listening on %s (env=%s)", addr, cfg.Env)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Auth-Token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
