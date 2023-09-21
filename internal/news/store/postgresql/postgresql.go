package postgresql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/news"
	"github.com/x-sports/internal/news/service"
)

var (
	errInvalidCommit   = errors.New("cannot do commit on non-transactional querier")
	errInvalidRollback = errors.New("cannot do rollback on non-transactional querier")
)

// store implements news/service.PGStore
type store struct {
	db *sqlx.DB
}

// storeClient implements news/service.PGStoreClient
type storeClient struct {
	q sqlx.Ext
}

// New creates a new store.
func New(db *sqlx.DB) (*store, error) {
	s := &store{
		db: db,
	}

	return s, nil
}

func (s *store) NewClient(useTx bool) (service.PGStoreClient, error) {
	var q sqlx.Ext

	// determine what object should be use as querier
	q = s.db
	if useTx {
		var err error
		q, err = s.db.Beginx()
		if err != nil {
			return nil, err
		}
	}

	return &storeClient{
		q: q,
	}, nil
}

func (sc *storeClient) Commit() error {
	if tx, ok := sc.q.(*sqlx.Tx); ok {
		return tx.Commit()
	}
	return errInvalidCommit
}

func (sc *storeClient) Rollback() error {
	if tx, ok := sc.q.(*sqlx.Tx); ok {
		return tx.Rollback()
	}
	return errInvalidRollback
}

// newsDB denotes a school data in the store.
type newsDB struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	GameID      int64      `db:"game_id"`
	GameNames   string     `db:"game_names"`
	GameIcons   string     `db:"game_icons"`
	Description string     `db:"description"`
	ImageNews   string     `db:"image_news"`
	Date        time.Time  `db:"date"`
	CreateTime  time.Time  `db:"create_time"`
	UpdateTime  *time.Time `db:"update_time"`
}

// format formats database struct into domain struct.
func (ndb *newsDB) format() news.News {
	n := news.News{
		ID:          ndb.ID,
		Title:       ndb.Title,
		GameID:      ndb.GameID,
		GameNames:   ndb.GameNames,
		GameIcons:   ndb.GameIcons,
		Description: ndb.Description,
		ImageNews:   ndb.ImageNews,
		Date:        ndb.Date,
		CreateTime:  ndb.CreateTime,
	}

	if ndb.UpdateTime != nil {
		n.UpdateTime = *ndb.UpdateTime
	}

	return n
}
