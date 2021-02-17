package internalhttp

import (
	"context"
	"fmt"
	"github.com/dimazusov/hw-test/advertising-banners/internal/domain"
	"net/http"

	"github.com/dimazusov/hw-test/advertising-banners/internal/config"
	"github.com/pkg/errors"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type server struct {
	app Application
	cfg *config.Config
	srv *http.Server
}

type Application interface {
	LogInfo(data interface{}) error
	LogError(data interface{}) error
	AddBannerToPlace(ctx context.Context, bannerID, placeID uint) (err error)
	DeleteBannerFromPlace(ctx context.Context, bannerID, placeID uint) (err error)
	AddEvent(ctx context.Context, event domain.Event) (err error)
	GetBannerForShow(ctx context.Context, placeID, socGroupID uint) (bannerID uint, err error)
}

func NewServer(cfg *config.Config, app Application) Server {
	return &server{
		cfg: cfg,
		app: app,
	}
}

func (m *server) Start(ctx context.Context) error {
	router := NewGinRouter(m.app)

	m.srv = &http.Server{}
	m.srv.Addr = m.cfg.Server.HTTP.Host + ":" + m.cfg.Server.HTTP.Port
	m.srv.Handler = router

	fmt.Println("listen: ", m.srv.Addr)

	err := m.srv.ListenAndServe()
	if err != nil {
		return errors.Wrap(err, "cannot listen and serve")
	}

	return nil
}

func (m *server) Stop(ctx context.Context) error {
	err := m.srv.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot shutdown server")
	}

	return nil
}
