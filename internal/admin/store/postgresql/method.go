package postgresql

import (
	"context"

	"github.com/x-sports/internal/admin"
)

func (sc *storeClient) GetUserByEmail(ctx context.Context, email string) (admin.Admin, error) {
	// query single row
	var adb adminDB
	err := sc.q.QueryRowx(queryGetUserByEmail, email).StructScan(&adb)
	if err != nil {
		return admin.Admin{}, err
	}

	return adb.format(), nil
}
