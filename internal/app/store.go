package app

import (
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Authenticate(username, password string) (User, bool, error) {
	row := s.db.QueryRow(
		"SELECT id, username, password, email, is_admin, balance, api_key FROM users WHERE username = ? AND password = ?",
		username,
		password,
	)

	user, err := scanUser(row)
	if err == sql.ErrNoRows {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, err
	}
	return user, true, nil
}

func (s *Store) ListUsers() ([]User, error) {
	rows, err := s.db.Query("SELECT id, username, password, email, is_admin, balance, api_key FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Store) GetUserByID(id int) (User, bool, error) {
	row := s.db.QueryRow(
		"SELECT id, username, password, email, is_admin, balance, api_key FROM users WHERE id = ?",
		id,
	)
	user, err := scanUser(row)
	if err == sql.ErrNoRows {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, err
	}
	return user, true, nil
}

func (s *Store) CreateUser(user User) (User, error) {
	res, err := s.db.Exec(
		"INSERT INTO users (username, password, email, is_admin, balance, api_key) VALUES (?, ?, ?, ?, ?, ?)",
		user.Username,
		user.Password,
		user.Email,
		boolToInt(user.IsAdmin),
		user.Balance,
		user.APIKey,
	)
	if err != nil {
		return User{}, err
	}
	id, _ := res.LastInsertId()
	user.ID = int(id)
	return user, nil
}

func (s *Store) ReplaceUser(id int, user User) (User, bool, error) {
	res, err := s.db.Exec(
		"UPDATE users SET username = ?, password = ?, email = ?, is_admin = ?, balance = ?, api_key = ? WHERE id = ?",
		user.Username,
		user.Password,
		user.Email,
		boolToInt(user.IsAdmin),
		user.Balance,
		user.APIKey,
		id,
	)
	if err != nil {
		return User{}, false, err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return User{}, false, nil
	}
	user.ID = id
	return user, true, nil
}

func (s *Store) PromoteUser(id int) (User, bool, error) {
	res, err := s.db.Exec("UPDATE users SET is_admin = 1 WHERE id = ?", id)
	if err != nil {
		return User{}, false, err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return User{}, false, nil
	}
	return s.GetUserByID(id)
}

func (s *Store) ResetPasswordByEmail(email, password string) (User, bool, error) {
	res, err := s.db.Exec("UPDATE users SET password = ? WHERE lower(email) = lower(?)", password, email)
	if err != nil {
		return User{}, false, err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return User{}, false, nil
	}

	row := s.db.QueryRow(
		"SELECT id, username, password, email, is_admin, balance, api_key FROM users WHERE lower(email) = lower(?)",
		email,
	)
	user, err := scanUser(row)
	if err != nil {
		return User{}, false, err
	}
	return user, true, nil
}

func (s *Store) ListProducts() ([]Product, error) {
	rows, err := s.db.Query("SELECT id, name, category, price, supplier, cost, internal_tag FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *Store) SearchProductsUnsafe(query string) ([]Product, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *Store) ListOrders() ([]Order, error) {
	rows, err := s.db.Query("SELECT id, user_id, product_id, quantity, address, status FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *Store) GetOrderByID(id int) (Order, bool, error) {
	row := s.db.QueryRow("SELECT id, user_id, product_id, quantity, address, status FROM orders WHERE id = ?", id)
	order, err := scanOrder(row)
	if err == sql.ErrNoRows {
		return Order{}, false, nil
	}
	if err != nil {
		return Order{}, false, err
	}
	return order, true, nil
}

func (s *Store) ListOrdersByUser(userID int) ([]Order, error) {
	rows, err := s.db.Query("SELECT id, user_id, product_id, quantity, address, status FROM orders WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *Store) CreateOrder(order Order) (Order, error) {
	res, err := s.db.Exec(
		"INSERT INTO orders (user_id, product_id, quantity, address, status) VALUES (?, ?, ?, ?, ?)",
		order.UserID,
		order.ProductID,
		order.Quantity,
		order.Address,
		order.Status,
	)
	if err != nil {
		return Order{}, err
	}
	id, _ := res.LastInsertId()
	order.ID = int(id)
	return order, nil
}

func scanUser(scanner interface{ Scan(dest ...any) error }) (User, error) {
	var user User
	var isAdmin int
	if err := scanner.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &isAdmin, &user.Balance, &user.APIKey); err != nil {
		return User{}, err
	}
	user.IsAdmin = isAdmin == 1
	return user, nil
}

func scanProduct(scanner interface{ Scan(dest ...any) error }) (Product, error) {
	var product Product
	if err := scanner.Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.Supplier, &product.Cost, &product.InternalTag); err != nil {
		return Product{}, err
	}
	return product, nil
}

func scanOrder(scanner interface{ Scan(dest ...any) error }) (Order, error) {
	var order Order
	if err := scanner.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Address, &order.Status); err != nil {
		return Order{}, err
	}
	return order, nil
}

func (s *Store) Ping() error {
	return s.db.Ping()
}

func (s *Store) DebugStats() (string, error) {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return "", err
	}
	return fmt.Sprintf("users=%d", count), nil
}
