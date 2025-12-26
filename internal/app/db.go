package app

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schemaSQL = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL,
	password TEXT NOT NULL,
	email TEXT NOT NULL,
	is_admin INTEGER NOT NULL DEFAULT 0,
	balance INTEGER NOT NULL DEFAULT 0,
	api_key TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	category TEXT NOT NULL,
	price REAL NOT NULL,
	supplier TEXT NOT NULL,
	cost REAL NOT NULL,
	internal_tag TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	product_id INTEGER NOT NULL,
	quantity INTEGER NOT NULL,
	address TEXT NOT NULL,
	status TEXT NOT NULL
);
`

func InitDatabase(path string) (*sql.DB, error) {
	if path == "" {
		path = "data/vulnapi.db"
	}

	if path != ":memory:" {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("create db dir: %w", err)
		}
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, fmt.Errorf("init schema: %w", err)
	}

	if err := seedData(db); err != nil {
		return nil, err
	}

	return db, nil
}

func seedData(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return nil
	}

	users := []User{
		{Username: "alice", Password: "password123", Email: "alice@harbor.example", IsAdmin: true, Balance: 9000, APIKey: "AKIA-ALICE-ADMIN"},
		{Username: "bob", Password: "bobpass", Email: "bob@harbor.example", IsAdmin: false, Balance: 120, APIKey: "AKIA-BOB-USER"},
		{Username: "charlie", Password: "charlie", Email: "charlie@harbor.example", IsAdmin: false, Balance: 60, APIKey: "AKIA-CHARLIE-USER"},
	}
	for _, u := range users {
		_, err := db.Exec(
			"INSERT INTO users (username, password, email, is_admin, balance, api_key) VALUES (?, ?, ?, ?, ?, ?)",
			u.Username, u.Password, u.Email, boolToInt(u.IsAdmin), u.Balance, u.APIKey,
		)
		if err != nil {
			return fmt.Errorf("seed users: %w", err)
		}
	}

	products := []Product{
		{Name: "Signal Dock", Category: "hardware", Price: 49.99, Supplier: "Acme", Cost: 12.00, InternalTag: "lab-only"},
		{Name: "Legacy Router", Category: "hardware", Price: 199.99, Supplier: "OldCo", Cost: 90.00, InternalTag: "clearance"},
		{Name: "Flow Monitor", Category: "services", Price: 29.99, Supplier: "SkyNet", Cost: 4.20, InternalTag: "beta"},
	}
	for _, p := range products {
		_, err := db.Exec(
			"INSERT INTO products (name, category, price, supplier, cost, internal_tag) VALUES (?, ?, ?, ?, ?, ?)",
			p.Name, p.Category, p.Price, p.Supplier, p.Cost, p.InternalTag,
		)
		if err != nil {
			return fmt.Errorf("seed products: %w", err)
		}
	}

	orders := []Order{
		{UserID: 2, ProductID: 1, Quantity: 2, Address: "12 Harbor St", Status: "processing"},
		{UserID: 3, ProductID: 2, Quantity: 1, Address: "88 Market Ave", Status: "shipped"},
	}
	for _, o := range orders {
		_, err := db.Exec(
			"INSERT INTO orders (user_id, product_id, quantity, address, status) VALUES (?, ?, ?, ?, ?)",
			o.UserID, o.ProductID, o.Quantity, o.Address, o.Status,
		)
		if err != nil {
			return fmt.Errorf("seed orders: %w", err)
		}
	}

	return nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
