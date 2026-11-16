package board

type Photo struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Advertisement struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Photos      []Photo `json:"photos"`
}

type Advertisements struct {
	Items []Advertisement `json:"items"`
}

type FavoriteItem struct {
	Advertisement Advertisement `json:"advertisement"`
}

type FavoriteList struct {
	Items []FavoriteItem `json:"items"`
	Total *int           `json:"total"`
}

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
