package postgresql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/team"
	"github.com/x-sports/internal/team/service"
)

var (
	errInvalidCommit   = errors.New("cannot do commit on non-transactional querier")
	errInvalidRollback = errors.New("cannot do rollback on non-transactional querier")
)

// store implements team/service.PGStore
type store struct {
	db *sqlx.DB
}

// storeClient implements team/service.PGStoreClient
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

// teamDB denotes a school data in the store.
type teamDB struct {
	ID         int64      `db:"id"`
	TeamNames  string     `db:"team_names"`
	GameID     int64      `db:"game_id"`
	GameNames  string     `db:"game_names"`
	GameIcons  string     `db:"game_icons"`
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
}

// format formats database struct into domain struct.
func (tdb *teamDB) format() team.Team {
	t := team.Team{
		ID:         tdb.ID,
		TeamNames:  tdb.TeamNames,
		GameID:     tdb.GameID,
		GameNames:  tdb.GameNames,
		GameIcons:  tdb.GameIcons,
		CreateTime: tdb.CreateTime,
	}

	if tdb.UpdateTime != nil {
		t.UpdateTime = *tdb.UpdateTime
	}

	return t
}
