package tg

import (
	"fmt"
	"github.com/ZergsLaw/tg-translate/internal/core"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type api struct {
	app    core.App
	bot    *tb.Bot
	logger *zap.Logger
}

var languages = [...]core.Lang{core.EN, core.DE, core.FR, core.ES, core.PT, core.IT, core.NL, core.PL, core.RU}
var langButtons [][]tb.ReplyButton

func convertLang(str string) core.Lang {
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
	default:
		return core.EN
	}
}

func New(app core.App) error {
	b, err := tb.NewBot(tb.Settings{
		URL:      "",
		Token:    "1085912034:AAFXApzG2RP-MHNunfmOaH7UGJlc1Uv0dOg",
		Updates:  0,
		Poller:   &tb.LongPoller{Timeout: 10 * time.Second},
		Reporter: nil,
		Client:   nil,
	})

	if err != nil {
		return fmt.Errorf("new bot: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}

	api := api{
		app:    app,
		bot:    b,
		logger: logger,
	}

	b.Handle("/start", api.start)

	revert := tb.ReplyButton{Text: "🔁"}
	b.Handle(&revert, api.revert)
	b.Handle(tb.OnText, api.translate)

	for i := range languages {
		lang := tb.ReplyButton{Text: languages[i].String()}
		langButtons = append(langButtons, []tb.ReplyButton{lang})
		b.Handle(&lang, api.changeLang(languages[i]))
	}

	for i := range languages {
		from := tb.ReplyButton{Text: fmt.Sprintf("from \n %s", languages[i])}
		to := tb.ReplyButton{Text: fmt.Sprintf("to \n %s", languages[i])}

		b.Handle(&from, api.changeFromLang())
		b.Handle(&to, api.changeToLang())
	}

	b.Start()

	return nil
}

func buttons(user core.User) [][]tb.ReplyButton {
	fromLang := tb.ReplyButton{Text: fmt.Sprintf("from \n %s", user.CurrentLangFrom)}
	revert := tb.ReplyButton{Text: "🔁"}
	toLang := tb.ReplyButton{Text: fmt.Sprintf("to \n %s", user.CurrentLangTo)}

	return [][]tb.ReplyButton{{fromLang, revert, toLang}}
}
