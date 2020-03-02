package main

import (
	"context"
	"fmt"
	"github.com/ZergsLaw/tg-translate/internal/core"
	"github.com/ZergsLaw/tg-translate/internal/modules/api/tg"
	"github.com/ZergsLaw/tg-translate/internal/modules/lang"
	"github.com/ZergsLaw/tg-translate/internal/modules/repo"
	"github.com/ZergsLaw/tg-translate/migrate"
	"log"
	"time"
)

const (
	timeoutConnect = time.Second * 5
)


func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutConnect)
	defer cancel()

	err := migrate.Run(ctx, "migrate", "up")
	if err != nil {
		panic(fmt.Errorf("migration: %w", err))
	}

	dbConn, err := repo.Connect(ctx)
	if err != nil {
		panic(fmt.Errorf("connect Repo: %w", err))
	}

	r := repo.New(dbConn)

	application := core.New(r, lang.New())

	log.Fatal(tg.New(application))
}
