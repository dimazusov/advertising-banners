package app

import (
	"context"
	"encoding/json"

	"github.com/dimazusov/hw-test/advertising-banners/internal/config"
	"github.com/dimazusov/hw-test/advertising-banners/internal/domain"
	"github.com/dimazusov/hw-test/advertising-banners/internal/kafka"
	"github.com/dimazusov/hw-test/advertising-banners/internal/pkg/ucb1"
)

type app struct {
	logger Logger
	rep    Repository
	cfg    *config.Config
}

type Logger interface {
	Debug(data interface{}) error
	Info(data interface{}) error
	Warn(data interface{}) error
	Error(data interface{}) error
	Close() error
}

type Repository interface {
	AddBannerToPlace(ctx context.Context, bannerID, placeID uint) (err error)
	DeleteBannerFromPlace(ctx context.Context, bannerID, placeID uint) (err error)
	AddEvent(ctx context.Context, event domain.Event) (err error)
	GetBannersStats(ctx context.Context, placeID, socGroupID uint) ([]domain.BannerStat, error)
	GetBannerIDsForPlace(ctx context.Context, placeID uint) ([]uint, error)
}

type App interface {
	LogInfo(interface{}) error
	LogError(interface{}) error
	AddBannerToPlace(ctx context.Context, bannerID, placeID uint) (err error)
	DeleteBannerFromPlace(ctx context.Context, bannerID, placeID uint) (err error)
	AddEvent(ctx context.Context, event domain.Event) (err error)
	GetBannerForShow(ctx context.Context, placeID, socGroupID uint) (bannerID uint, err error)
}

func New(logger Logger, repository Repository, cfg *config.Config) App {
	return &app{
		logger: logger,
		rep:    repository,
		cfg:    cfg,
	}
}

func (m *app) LogInfo(data interface{}) error {
	return m.logger.Info(data)
}

func (m *app) LogError(data interface{}) error {
	return m.logger.Info(data)
}

func (m *app) AddBannerToPlace(ctx context.Context, bannerID, placeID uint) (err error) {
	return m.rep.AddBannerToPlace(ctx, bannerID, placeID)
}

func (m *app) DeleteBannerFromPlace(ctx context.Context, bannerID, placeID uint) (err error) {
	return m.rep.DeleteBannerFromPlace(ctx, bannerID, placeID)
}

func (m *app) AddEvent(ctx context.Context, event domain.Event) (err error) {
	producer, err := kafka.NewProducer("event", m.cfg)
	if err != nil {
		return err
	}

	b, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	err = producer.WriteMessages(context.Background(), b)
	if err != nil {
		return err
	}

	return m.rep.AddEvent(ctx, event)
}

func (m *app) GetBannerForShow(ctx context.Context, placeID, socGroupID uint) (bannerID uint, err error) {
	bannerIDs, err := m.rep.GetBannerIDsForPlace(ctx, placeID)
	if err != nil {
		return 0, err
	}

	bannerStatistics, err := m.rep.GetBannersStats(ctx, placeID, socGroupID)
	if err != nil {
		return 0, err
	}

	for _, bannerID := range bannerIDs {
		isFinded := false
		for _, stat := range bannerStatistics {
			if bannerID == stat.ID {
				isFinded = true
			}
		}
		if !isFinded {
			return bannerID, nil
		}
	}

	alg := ucb1.New()
	for _, stat := range bannerStatistics {
		alg.Add(stat.ID, stat.Count, stat.Total)
	}

	return alg.GetBest(), nil
}
