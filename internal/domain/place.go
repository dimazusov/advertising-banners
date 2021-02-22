package domain

// Place for banner
type Place struct {
	ID          uint   `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
}
