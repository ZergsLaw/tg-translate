package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ZergsLaw/tg-translate/internal/core"
)

func (repo *Repo) CreateUser(ctx context.Context, newUser core.User) (core.UserID, error) {
	const query = `INSERT INTO users (telegram_id, current_lang_from, current_lang_to) VALUES ($1,$2,$3)`

	_, err := repo.db.ExecContext(ctx, query, newUser.TelegramID, newUser.CurrentLangFrom.String(), newUser.CurrentLangTo.String())
	if err != nil {
		return 0, fmt.Errorf("insert new user: %w", err)
	}

	user, err := repo.UserByTelegramID(ctx, newUser.TelegramID)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (repo *Repo) UserByID(ctx context.Context, userID core.UserID) (*core.User, error) {
	const query = `SELECT * FROM users WHERE id = $1`

	u := &user{}
	err := repo.db.QueryRowContext(ctx, query, userID).Scan(
		&u.ID,
		&u.TelegramID,
		&u.CurrentLangFrom,
		&u.CurrentLangTo,
		&u.CurrentLangState,
		&u.CreatedAt,
		&u.LastActionTime,
	)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, core.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return u.Convert(), nil
}

func (repo *Repo) UserByTelegramID(ctx context.Context, tgID core.TelegramID) (*core.User, error) {
	const query = `SELECT * FROM users WHERE telegram_id = $1`

	u := &user{}
	err := repo.db.QueryRowContext(ctx, query, tgID).Scan(
		&u.ID,
		&u.TelegramID,
		&u.CurrentLangFrom,
		&u.CurrentLangTo,
		&u.CurrentLangState,
		&u.CreatedAt,
		&u.LastActionTime,
	)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, core.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("get user by telegram id: %w", err)
	}

	return u.Convert(), nil
}

func (repo *Repo) SetCurrentLang(ctx context.Context, id core.UserID, from, to core.Lang) error {
	const query = `UPDATE users SET current_lang_from = $1, current_lang_to = $2, current_lang_state = null WHERE id = $3`

	_, err := repo.db.ExecContext(ctx, query, from.String(), to.String(), id)
	if err != nil {
		return fmt.Errorf("update lang: %w", err)
	}

	return nil
}

func (repo *Repo) SetCurrentLangState(ctx context.Context, id core.UserID, state core.CurrentLangStateEdit) error {
	const query = `UPDATE users SET current_lang_state = $1 WHERE id = $2`

	_, err := repo.db.ExecContext(ctx, query, state.String(), id)
	if err != nil {
		return fmt.Errorf("update lang: %w", err)
	}

	return nil
}

func (repo *Repo) UpdateLastActionTime(ctx context.Context, id core.UserID) error {
	const query = `UPDATE users SET last_action_time = now() WHERE id = $1`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("update last action time: %w", err)
	}

	return nil
}
