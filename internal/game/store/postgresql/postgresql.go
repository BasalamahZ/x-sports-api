package postgresql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/game"
	"github.com/x-sports/internal/game/service"
)

var (
	errInvalidCommit   = errors.New("cannot do commit on non-transactional querier")
	errInvalidRollback = errors.New("cannot do rollback on non-transactional querier")
)

// store implements game/service.PGStore
type store struct {
	db *sqlx.DB
}

// storeClient implements game/service.PGStoreClient
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

// gameDB denotes a school data in the store.
type gameDB struct {
	ID         int64      `db:"id"`
	GameNames  string     `db:"game_names"`
	GameIcons  string     `db:"game_icons"`
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
}

// format formats database struct into domain struct.
func (gdb *gameDB) format() game.Game {
	g := game.Game{
		ID:         gdb.ID,
		GameNames:  gdb.GameNames,
		GameIcons:  gdb.GameIcons,
		CreateTime: gdb.CreateTime,
	}

	if gdb.UpdateTime != nil {
		g.UpdateTime = *gdb.UpdateTime
	}

	return g
}
