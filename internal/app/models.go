package app

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
	Balance  int    `json:"balance"`
	APIKey   string `json:"apiKey"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Supplier    string  `json:"supplier"`
	Cost        float64 `json:"cost"`
	InternalTag string  `json:"internalTag"`
}

type Order struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	ProductID int    `json:"productId"`
	Quantity  int    `json:"quantity"`
	Address   string `json:"address"`
	Status    string `json:"status"`
}

type AuthInfo struct {
	UserID  int    `json:"userId"`
	IsAdmin bool   `json:"isAdmin"`
	Token   string `json:"token"`
}
