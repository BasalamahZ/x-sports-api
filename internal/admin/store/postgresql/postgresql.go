package postgresql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/admin"
	"github.com/x-sports/internal/admin/service"
)

var (
	errInvalidCommit   = errors.New("cannot do commit on non-transactional querier")
	errInvalidRollback = errors.New("cannot do rollback on non-transactional querier")
)

// store implements admin/service.PGStore
type store struct {
	db *sqlx.DB
}

// storeClient implements admin/service.PGStoreClient
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

// adminDB denotes a school data in the store.
type adminDB struct {
	ID         int64      `db:"id"`
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
}

// format formats database struct into domain struct.
func (adb *adminDB) format() admin.Admin {
	var updateTime time.Time
	if adb.UpdateTime != nil {
		updateTime = *adb.UpdateTime
	}

	a := admin.Admin{
		ID:         adb.ID,
		Email:      adb.Email,
		Password:   adb.Password,
		CreateTime: adb.CreateTime,
		UpdateTime: updateTime,
	}

	return a
}
