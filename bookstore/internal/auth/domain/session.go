package domain

import "time"

type Session struct {
	EMail   string    `json:"email"`
	Expires time.Time `json:"expires"`
}
