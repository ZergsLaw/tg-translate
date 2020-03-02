package repo_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ZergsLaw/tg-translate/internal/core"
	"github.com/ZergsLaw/tg-translate/internal/modules/repo"
	"github.com/ZergsLaw/tg-translate/migrate"
)

var (
	Repo core.Repo
)

const (
	timeoutConnect = time.Second * 100000
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutConnect)
	defer cancel()

	resetDB := func() {
		err := migrate.Run(ctx, "../../../migrate", "reset")
		if err != nil {
			panic(fmt.Errorf("migration: %w", err))
		}
	}
	// For convenient cleaning DB.
	resetDB()

	err := migrate.Run(ctx, "../../../migrate", "up")
	if err != nil {
		panic(fmt.Errorf("migration: %w", err))
	}

	defer resetDB()

	dbConn, err := repo.Connect(ctx)
	if err != nil {
		panic(fmt.Errorf("connect Repo: %w", err))
	}

	Repo = repo.New(dbConn)

	os.Exit(m.Run())
}
