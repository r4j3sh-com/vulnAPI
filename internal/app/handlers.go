package app

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func loginHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := decodeJSON(r.Body, &creds); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		user, ok, err := store.Authenticate(creds.Username, creds.Password)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "login failed")
			return
		}
		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		token := IssueToken(user)
		writeJSON(w, http.StatusOK, map[string]any{
			"token": token,
			"user":  user,
		})
	}
}

func passwordResetHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var payload struct {
			Email       string `json:"email"`
			NewPassword string `json:"newPassword"`
		}

		if err := decodeJSON(r.Body, &payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		user, ok, err := store.ResetPasswordByEmail(payload.Email, payload.NewPassword)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "reset failed")
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"message": "password reset",
			"user":    user,
		})
	}
}

func dashboardHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		auth, ok := requireAuth(w, r)
		if !ok {
			return
		}

		user, found, err := store.GetUserByID(auth.UserID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load user")
			return
		}
		if !found {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}

		writeJSON(w, http.StatusOK, user)
	}
}

func usersHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			users, err := store.ListUsers()
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to list users")
				return
			}
			writeJSON(w, http.StatusOK, users)
		case http.MethodPost:
			var user User
			if err := decodeJSON(r.Body, &user); err != nil {
				writeError(w, http.StatusBadRequest, "invalid json")
				return
			}

			created, err := store.CreateUser(user)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "create failed")
				return
			}
			writeJSON(w, http.StatusCreated, created)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func userHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/users/")
		parts := strings.Split(path, "/")
		if len(parts) == 0 || parts[0] == "" {
			writeError(w, http.StatusBadRequest, "missing id")
			return
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		if len(parts) > 1 && parts[1] == "orders" {
			orders, err := store.ListOrdersByUser(id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to list orders")
				return
			}
			writeJSON(w, http.StatusOK, orders)
			return
		}

		switch r.Method {
		case http.MethodGet:
			user, ok, err := store.GetUserByID(id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to get user")
				return
			}
			if !ok {
				writeError(w, http.StatusNotFound, "user not found")
				return
			}
			writeJSON(w, http.StatusOK, user)
		case http.MethodPut:
			var payload User
			if err := decodeJSON(r.Body, &payload); err != nil {
				writeError(w, http.StatusBadRequest, "invalid json")
				return
			}
			updated, ok, err := store.ReplaceUser(id, payload)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "update failed")
				return
			}
			if !ok {
				writeError(w, http.StatusNotFound, "user not found")
				return
			}
			writeJSON(w, http.StatusOK, updated)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func ordersHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orders, err := store.ListOrders()
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to list orders")
				return
			}
			writeJSON(w, http.StatusOK, orders)
		case http.MethodPost:
			if _, ok := requireAuth(w, r); !ok {
				return
			}

			var payload Order
			if err := decodeJSON(r.Body, &payload); err != nil {
				writeError(w, http.StatusBadRequest, "invalid json")
				return
			}

			created, err := store.CreateOrder(payload)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "create failed")
				return
			}
			writeJSON(w, http.StatusCreated, created)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func orderHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requireAuth(w, r); !ok {
			return
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/orders/"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		switch r.Method {
		case http.MethodGet:
			order, ok, err := store.GetOrderByID(id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to get order")
				return
			}
			if !ok {
				writeError(w, http.StatusNotFound, "order not found")
				return
			}
			writeJSON(w, http.StatusOK, order)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func productsHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			products, err := store.ListProducts()
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to list products")
				return
			}
			writeJSON(w, http.StatusOK, products)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func searchHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		term := r.URL.Query().Get("q")
		query := "SELECT * FROM products WHERE name LIKE '%" + term + "%'"
		results, err := store.SearchProductsUnsafe(query)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "search failed")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"query":   query,
			"results": results,
		})
	}
}

func reportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		sizeStr := r.URL.Query().Get("size")
		size, _ := strconv.Atoi(sizeStr)
		if size <= 0 {
			size = 1024
		}

		payload := strings.Repeat("A", size)
		writeJSON(w, http.StatusOK, map[string]any{
			"generated": size,
			"data":      payload,
		})
	}
}

func adminStatsHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if _, ok := requireAuth(w, r); !ok {
			return
		}

		users, err := store.ListUsers()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load stats")
			return
		}
		orders, err := store.ListOrders()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load stats")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"users":  len(users),
			"orders": len(orders),
			"uptime": time.Now().Unix(),
		})
	}
}

func promoteHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if _, ok := requireAuth(w, r); !ok {
			return
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/ops/roles/promote/"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		user, ok, err := store.PromoteUser(id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "promote failed")
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}

		writeJSON(w, http.StatusOK, user)
	}
}

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		host := r.URL.Query().Get("host")
		if host == "" {
			host = "127.0.0.1"
		}

		cmd := exec.Command("sh", "-c", "ping -c 1 "+host)
		output, err := cmd.CombinedOutput()
		if err != nil {
			writeJSON(w, http.StatusOK, map[string]any{
				"host":   host,
				"error":  err.Error(),
				"output": string(output),
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"host":   host,
			"output": string(output),
		})
	}
}

func fileReadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := r.URL.Query().Get("path")
		if path == "" {
			path = "/etc/hosts"
		}

		data, err := os.ReadFile(path)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func v1UsersHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		users, err := store.ListUsers()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to list users")
			return
		}
		writeJSON(w, http.StatusOK, users)
	}
}

func v2UsersHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		users, err := store.ListUsers()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to list users")
			return
		}
		for i := range users {
			users[i].Password = ""
			users[i].APIKey = ""
		}
		writeJSON(w, http.StatusOK, users)
	}
}

func debugConfigHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"env":         cfg.Env,
			"debug":       cfg.Debug,
			"tokenSecret": cfg.TokenSecret,
			"dbPath":      cfg.DBPath,
		})
	}
}

func requireAuth(w http.ResponseWriter, r *http.Request) (AuthInfo, bool) {
	auth, err := AuthFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing or invalid token")
		return AuthInfo{}, false
	}
	return auth, true
}

func decodeJSON(body io.Reader, dest any) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(dest)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg})
}
