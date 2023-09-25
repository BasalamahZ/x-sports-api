package postgresql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/match"
	"github.com/x-sports/internal/match/service"
)

var (
	errInvalidCommit   = errors.New("cannot do commit on non-transactional querier")
	errInvalidRollback = errors.New("cannot do rollback on non-transactional querier")
)

// store implements match/service.PGStore
type store struct {
	db *sqlx.DB
}

// storeClient implements match/service.PGStoreClient
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

// matchDB denotes a school data in the store.
type matchDB struct {
	ID              int64        `db:"id"`
	BlockChainID    uint64       `db:"blockchain_id"`
	TournamentNames string       `db:"tournament_names"`
	GameID          int64        `db:"game_id"`
	GameNames       string       `db:"game_names"`
	GameIcons       string       `db:"game_icons"`
	TeamAID         int64        `db:"team_a_id"`
	TeamANames      string       `db:"team_a_names"`
	TeamAIcons      string       `db:"team_a_icons"`
	TeamAOdds       float32      `db:"team_a_odds"`
	TeamBID         int64        `db:"team_b_id"`
	TeamBNames      string       `db:"team_b_names"`
	TeamBIcons      string       `db:"team_b_icons"`
	TeamBOdds       float32      `db:"team_b_odds"`
	Date            time.Time    `db:"date"`
	MatchLink       string       `db:"match_link"`
	Status          match.Status `db:"status"`
	Winner          int64        `db:"winner"`
	CreateTime      time.Time    `db:"create_time"`
	UpdateTime      *time.Time   `db:"update_time"`
}

// format formats database struct into domain struct.
func (mdb *matchDB) format() match.Match {
	m := match.Match{
		ID:              mdb.ID,
		BlockChainID:    mdb.BlockChainID,
		TournamentNames: mdb.TournamentNames,
		GameID:          mdb.GameID,
		GameNames:       mdb.GameNames,
		GameIcons:       mdb.GameIcons,
		TeamAID:         mdb.TeamAID,
		TeamANames:      mdb.TeamANames,
		TeamAIcons:      mdb.TeamAIcons,
		TeamAOdds:       mdb.TeamAOdds,
		TeamBID:         mdb.TeamBID,
		TeamBNames:      mdb.TeamBNames,
		TeamBIcons:      mdb.TeamBIcons,
		TeamBOdds:       mdb.TeamBOdds,
		Date:            mdb.Date,
		Status:          mdb.Status,
		MatchLink:       mdb.MatchLink,
		Winner:          mdb.Winner,
		CreateTime:      mdb.CreateTime,
	}

	if mdb.UpdateTime != nil {
		m.UpdateTime = *mdb.UpdateTime
	}

	return m
}
