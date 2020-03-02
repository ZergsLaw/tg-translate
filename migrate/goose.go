package migrate

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ZergsLaw/tg-translate/internal/modules/repo"
	"github.com/pressly/goose"
)

// nolint:gochecknoglobals
var (
	gooseMu sync.Mutex
)

// Run executes goose command. It also enforce "fix" after "create".
func Run(ctx context.Context, dir string, command string, options ...repo.Option) error {
	gooseMu.Lock()
	defer gooseMu.Unlock()

	dbConn, err := repo.Connect(ctx, options...)
	if err != nil {
		return err
	}

	cmdArgs := strings.Fields(command)
	cmd, args := cmdArgs[0], cmdArgs[1:]
	err = goose.Run(cmd, dbConn, dir, args...)
	if err == nil && cmd == "create" {
		err = goose.Run("fix", dbConn, dir)
	}
	if err != nil {
		return fmt.Errorf("goose.Run %q: %w", command, err)
	}

	return dbConn.Close()
}
