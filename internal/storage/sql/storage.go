//nolint:golint
package sql

import (
	"context"
	"github.com/dimazusov/hw-test/advertising-banners/internal/domain"
	"github.com/dimazusov/hw-test/advertising-banners/internal/pkg/apperror"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type postgresStorage struct {
	conn       *sqlx.DB
	dsn        string
	driverName string
}

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	AddBannerToPlace(ctx context.Context, bannerID, placeID uint) (err error)
	DeleteBannerFromPlace(ctx context.Context, bannerID, placeID uint) (err error)
	AddEvent(ctx context.Context, event domain.Event) (err error)
	GetBannersStats(ctx context.Context, placeID, socGroupID uint) (stats []domain.BannerStat, err error)
	GetBannerIDsForPlace(ctx context.Context, placeID uint) ([]uint, error)
}

func New(driverName, dsn string) (Storage, error) {
	conn, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect")
	}

	return &postgresStorage{
		conn:       conn,
		dsn:        dsn,
		driverName: driverName,
	}, nil
}

func (m *postgresStorage) Connect(ctx context.Context) error {
	err := m.conn.PingContext(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot connect")
	}

	return nil
}

func (m *postgresStorage) Close() error {
	err := m.conn.Close()
	if err != nil {
		return errors.Wrap(err, "cannot close")
	}

	return nil
}

func (m *postgresStorage) AddBannerToPlace(ctx context.Context, bannerID, placeID uint) (err error) {
	_, err = m.conn.ExecContext(ctx, "insert into banner_place(banner_id, place_id) VALUES($1, $2)", bannerID, placeID)
	if err != nil {
		return errors.Wrap(err, "cannot insert to banner_place")
	}

	return nil
}

func (m *postgresStorage) DeleteBannerFromPlace(ctx context.Context, bannerID, placeID uint) (err error) {
	res, err := m.conn.ExecContext(ctx, "delete from banner_place where banner_id=$1 and place_id=$2", bannerID, placeID)
	if err != nil {
		return errors.Wrap(err, "cannot delete from banner_place")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get banner_place deleted rows")
	}

	if rowsAffected == 0 {
		return errors.Wrap(apperror.ErrNotFound, "banner for place not found")
	}

	return nil
}

func (m *postgresStorage) AddEvent(ctx context.Context, e domain.Event) (err error) {
	query := "insert into event(type, place_id, banner_id, soc_group_id, time) VALUES($1, $2, $3, $4, $5)"

	_, err = m.conn.ExecContext(ctx, query, e.Type, e.PlaceID, e.BannerID, e.SocGroupID, e.Time)
	if err != nil {
		return errors.Wrap(err, "cannot insert to banner_place")
	}

	return nil
}

func (m *postgresStorage) GetBannerIDsForPlace(ctx context.Context, placeID uint) ([]uint, error) {
	query := "select banner_id from banner_place where place_id = $1"
	rows, err := m.conn.QueryContext(ctx, query, placeID)
	if err != nil {
		return nil, err
	}

	bannerIDs := []uint{}
	for rows.Next() {
		var bannerID uint
		err := rows.Scan(&bannerID)
		if err != nil {
			return nil, err
		}
		bannerIDs = append(bannerIDs, bannerID)
	}

	return bannerIDs, nil
}

func (m *postgresStorage) GetBannersStats(ctx context.Context, placeID, socGroupID uint) (stats []domain.BannerStat, err error) {
	var total uint
	query := "select count(*) from event where place_id = $1 and soc_group_id = $2"
	err = m.conn.QueryRowxContext(ctx, query, placeID, socGroupID).Scan(&total)
	if err != nil {
		return nil, err
	}

	query = "select count(id), banner_id from event where place_id = $1 and soc_group_id = $2 GROUP BY banner_id"
	rows, err := m.conn.QueryContext(ctx, query, placeID, socGroupID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		stat := domain.BannerStat{}
		err := rows.Scan(&stat.Count, &stat.ID)
		if err != nil {
			return nil, err
		}

		stat.Total = total
		stats = append(stats, stat)
	}

	return stats, nil
}

