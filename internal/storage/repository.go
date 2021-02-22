package storage

import (
	"context"

	"github.com/dimazusov/hw-test/advertising-banners/internal/config"
	"github.com/dimazusov/hw-test/advertising-banners/internal/domain"
	sqlstorage "github.com/dimazusov/hw-test/advertising-banners/internal/storage/sql"
	"github.com/pkg/errors"
)

var ErrRepositoryTypeNotExists = errors.New("repository type not exists")

const RepTypePostgres = "postgres"

type Repository interface {
	AddBannerToPlace(ctx context.Context, bannerID, placeID uint) (err error)
	DeleteBannerFromPlace(ctx context.Context, bannerID, placeID uint) (err error)
	AddEvent(ctx context.Context, event domain.Event) (err error)
	GetBannersStats(ctx context.Context, placeID, socGroupID uint) (stats []domain.BannerStat, err error)
	GetBannerIDsForPlace(ctx context.Context, placeID uint) ([]uint, error)
}

func NewRepository(cfg *config.Config) (Repository, error) {
	var storage interface{}
	var err error

	switch cfg.Repository.Type {
	case RepTypePostgres:
		storage, err = sqlstorage.New(cfg.DB.Postgres.Dialect, cfg.DB.Postgres.Dsn)
	default:
		err = ErrRepositoryTypeNotExists
	}

	if err != nil {
		return nil, errors.Wrap(err, "cannot create repository")
	}

	return storage.(Repository), nil
}
