package repo

import (
	"database/sql"
	"fmt"
	"github.com/ZergsLaw/tg-translate/internal/core"
	"time"
)

type user struct {
	ID               int
	TelegramID       int
	CurrentLangFrom  string
	CurrentLangTo    string
	CurrentLangState sql.NullString
	CreatedAt        time.Time
	LastActionTime   time.Time
}

func (u *user) Convert() *core.User {
	return &core.User{
		ID:               core.UserID(u.ID),
		TelegramID:       core.TelegramID(u.TelegramID),
		CurrentLangFrom:  parseLang(u.CurrentLangFrom),
		CurrentLangTo:    parseLang(u.CurrentLangTo),
		CurrentLangState: parseLangEditState(u.CurrentLangState),
		CreatedAt:        u.CreatedAt,
		LastActionTime:   u.LastActionTime,
	}
}

func parseLangEditState(str sql.NullString) core.CurrentLangStateEdit {
	switch str.String {
	case "from":
		return core.From
	case "to":
		return core.To
	case "":
		return core.None
	}

	panic(fmt.Sprintf("unknown lang %s", str.String))
}

func parseLang(str string) core.Lang {
	switch str {
	case `German`:
		return core.DE
	case `French`:
		return core.FR
	case `Spanish`:
		return core.ES
	case `Portuguese`:
		return core.PT
	case `Italian`:
		return core.IT
	case `Dutch`:
		return core.DE
	case `Polish`:
		return core.PL
	case `Russian`:
		return core.RU
	case `English`:
		return core.EN
	}

	panic(fmt.Sprintf("unknown lang %s", str))
}
