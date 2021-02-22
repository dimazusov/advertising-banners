package internalhttp

//AddBannerToPlace(ctx context.Context, bannerID, placeID uint) (newID uint, err error)
//DeleteBannerFromPlace(ctx context.Context, bannerID, placeID uint) (err error)
//AddEvent(ctx context.Context, placeID, bannerID, socGroupID uint) (err error)
//GetBannerForShow(ctx context.Context, placeID, socGroupID uint) (bannerID uint, err error)

type addBannerToPlaceParams struct {
	BannerID uint `json:"bannerId"`
	PlaceID  uint `json:"placeId"`
}

type deleteBannerToPlaceParams addBannerToPlaceParams

type createEventParams struct {
	Type       string `json:"type"`
	PlaceID    uint   `json:"placeId"`
	BannerID   uint   `json:"bannerId"`
	SocGroupID uint   `json:"socGroupId"`
}

type bannerParams struct {
	PlaceID    uint `json:"placeId"`
	SocGroupID uint `json:"socGroupId"`
}
