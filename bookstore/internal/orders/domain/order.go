package domain

type Order struct {
	User     string `json:"user"`
	Contents []Book `json:"contents"`
}
