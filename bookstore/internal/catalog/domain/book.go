package domain

type Book struct {
	ID       string `json:"id"` // created by bbolt's `NextSequence()`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
