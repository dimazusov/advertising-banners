package domain

type Banner struct {
	ID          uint   `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
}

type BannerStat struct {
	ID    uint `json:"id"`
	Count uint `json:"count"`
	Total uint `json:"total"`
}
