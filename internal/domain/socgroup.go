package domain

// SocGroup for banner
type SocGroup struct {
	ID          uint   `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
}
