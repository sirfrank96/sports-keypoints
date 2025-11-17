package controller

import (
	"context"
	"fmt"

	db "github.com/sirfrank96/go-server/db"
)

func verifyUserExists(ctx context.Context, dbmgr *db.DbManager, userId string) (*db.User, error) {
	user, err := dbmgr.ReadUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("could not find user %s: %w", userId, err)
	}
	return user, nil
}
