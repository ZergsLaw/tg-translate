package core

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type (
	Repo interface {
		CreateUser(context.Context, User) (UserID, error)
		UserByID(context.Context, UserID) (*User, error)
		UserByTelegramID(context.Context, TelegramID) (*User, error)
		SetCurrentLang(ctx context.Context, id UserID, from, to Lang) error
		SetCurrentLangState(context.Context, UserID, CurrentLangStateEdit) error
		UpdateLastActionTime(context.Context, UserID) error
	}

	Translator interface {
		Translate(ctx context.Context, text string, from, to Lang) (string, error)
	}

	App interface {
		CreateOrGetUser(ctx context.Context, tgID TelegramID, userLang Lang) (*User, error)
		Translate(ctx context.Context, user User, text string) (string, error)
		RevertLang(ctx context.Context, user User) error
		SetCurrentLangFrom(ctx context.Context, user User, lang Lang) error
		SetCurrentLangTo(ctx context.Context, user User, lang Lang) error
		SetCurrentLangState(context.Context, User, CurrentLangStateEdit) error
	}

	Lang                 uint8
	CurrentLangStateEdit uint8

	UserID     int
	TelegramID int

	User struct {
		ID               UserID
		TelegramID       TelegramID
		CurrentLangFrom  Lang
		CurrentLangTo    Lang
		CurrentLangState CurrentLangStateEdit
		CreatedAt        time.Time
		LastActionTime   time.Time
	}
)

const (
	EN Lang = iota + 1
	DE
	FR
	ES
	PT
	IT
	NL
	PL
	RU
)

func (l Lang) String() string {
	switch l {
	case DE:
		return `German`
	case FR:
		return `French`
	case ES:
		return `Spanish`
	case PT:
		return `Portuguese`
	case IT:
		return `Italian`
	case NL:
		return `Dutch`
	case PL:
		return `Polish`
	case RU:
		return `Russian`
	case EN:
		return `English`
	}

	panic(fmt.Sprintf("unknown language %d", l))
}

func (state CurrentLangStateEdit) String() string {
	switch state {
	case From:
		return `from`
	case To:
		return `to`
	}

	panic(fmt.Sprintf("unknown state %d", state))
}

const (
	None CurrentLangStateEdit = iota + 1
	From
	To
)

type app struct {
	r Repo
	t Translator
}

func (a *app) SetCurrentLangState(ctx context.Context, user User, state CurrentLangStateEdit) error {
	return a.r.SetCurrentLangState(ctx, user.ID, state)
}

var (
	ErrNotFound = errors.New("not found")
)

func (a *app) CreateOrGetUser(ctx context.Context, tgID TelegramID, userLang Lang) (*User, error) {
	user, err := a.r.UserByTelegramID(ctx, tgID)
	switch {
	case err == nil:
		err = a.r.UpdateLastActionTime(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		return user, nil
	case errors.Is(err, ErrNotFound):
		user = &User{
			TelegramID:      tgID,
			CurrentLangFrom: userLang,
			CurrentLangTo:   alternativeLang(userLang),
		}

		userID, err := a.r.CreateUser(ctx, *user)
		if err != nil {
			return nil, err
		}
		user.ID = userID

		return user, nil
	default:
		return nil, err
	}
}

func (a *app) Translate(ctx context.Context, user User, text string) (string, error) {
	result, err := a.t.Translate(ctx, text, user.CurrentLangFrom, user.CurrentLangTo)
	if err != nil {
		return "", nil
	}

	err = a.r.UpdateLastActionTime(ctx, user.ID)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (a *app) RevertLang(ctx context.Context, user User) error {
	err := a.r.SetCurrentLang(ctx, user.ID, user.CurrentLangTo, user.CurrentLangFrom)
	if err != nil {
		return err
	}

	return a.r.UpdateLastActionTime(ctx, user.ID)
}

func (a *app) SetCurrentLangFrom(ctx context.Context, user User, lang Lang) error {
	err := a.r.SetCurrentLang(ctx, user.ID, lang, user.CurrentLangTo)
	if err != nil {
		return err
	}

	return a.r.UpdateLastActionTime(ctx, user.ID)
}

func (a *app) SetCurrentLangTo(ctx context.Context, user User, lang Lang) error {
	err := a.r.SetCurrentLang(ctx, user.ID, user.CurrentLangFrom, lang)
	if err != nil {
		return err
	}

	return a.r.UpdateLastActionTime(ctx, user.ID)
}

func New(r Repo, t Translator) App {
	return &app{
		r: r,
		t: t,
	}
}

// Except EN.
var langArray = [...]Lang{DE, FR, ES, PT, IT, NL, PL, RU}

func alternativeLang(userLang Lang) Lang {
	if userLang == EN {
		randInt := rand.Intn(len(langArray)) - 1
		return langArray[randInt]
	}

	return EN
}
