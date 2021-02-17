package domain

const (
	Click = "click"
	View  = "view"
)

type Event struct {
	Type       string `json:"type" db:"type"`
	PlaceID    uint   `json:"placeId" db:"place_id"`
	BannerID   uint   `json:"bannerId" db:"banner_id"`
	SocGroupID uint   `json:"socGroupId" db:"soc_group_id"`
	Time       uint   `json:"time" db:"time"`
}
