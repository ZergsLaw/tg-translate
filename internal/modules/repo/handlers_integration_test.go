package repo_test

import (
	"context"
	"testing"

	"github.com/ZergsLaw/tg-translate/internal/core"
	"github.com/stretchr/testify/require"
)

var (
	ctx = context.Background()
)

func TestRepoSmoke(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	user := core.User{
		TelegramID:       1,
		CurrentLangFrom:  core.EN,
		CurrentLangTo:    core.RU,
		CurrentLangState: core.None,
	}

	userID, err := Repo.CreateUser(ctx, user)
	assert.NoError(err)
	assert.NotNil(userID)
	user.ID = userID

	userFromRepo, err := Repo.UserByID(ctx, userID)
	assert.NoError(err)
	user.CreatedAt = userFromRepo.CreatedAt
	user.LastActionTime = userFromRepo.LastActionTime
	assert.Equal(&user, userFromRepo)

	err = Repo.SetCurrentLangState(ctx, user.ID, core.From)
	assert.NoError(err)
	user.CurrentLangState = core.From

	userFromRepo, err = Repo.UserByTelegramID(ctx, user.TelegramID)
	assert.NoError(err)
	assert.Equal(&user, userFromRepo)

	err = Repo.SetCurrentLang(ctx, user.ID, core.FR, core.PT)
	assert.NoError(err)
	user.CurrentLangFrom, user.CurrentLangTo = core.FR, core.PT
	user.CurrentLangState = core.None

	err = Repo.UpdateLastActionTime(ctx, user.ID)
	assert.NoError(err)

	userFromRepo, err = Repo.UserByTelegramID(ctx, user.TelegramID)
	assert.NoError(err)

	if !userFromRepo.LastActionTime.After(user.LastActionTime) {
		assert.Fail("last action time not correct")
	}
	user.LastActionTime = userFromRepo.LastActionTime
	assert.Equal(&user, userFromRepo)
}
