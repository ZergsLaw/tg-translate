package tg

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/ZergsLaw/tg-translate/internal/core"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

const timeout = time.Second * 3

func (api *api) start(m *tb.Message) {
	api.execReq(m.Sender, func(ctx context.Context, user *core.User) (message string, err error) {
		const welcomeMessage = "Hi! To translate the text, just send a message. " +
			"To change the translation languages, select them using the keyboard keys. :)"

		return welcomeMessage, nil
	})
}

func (api *api) revert(m *tb.Message) {
	api.execReq(m.Sender, func(ctx context.Context, user *core.User) (message string, err error) {
		err = api.app.RevertLang(ctx, *user)
		if err != nil {
			return "", err
		}
		user.CurrentLangFrom, user.CurrentLangTo = user.CurrentLangTo, user.CurrentLangFrom
		const successMessage = "Revert you language!"

		return successMessage, nil
	})
}

func (api *api) translate(m *tb.Message) {
	api.execReq(m.Sender, func(ctx context.Context, user *core.User) (message string, err error) {
		res, err := api.app.Translate(ctx, *user, m.Text)
		if err != nil {
			return "", err
		}

		return res, err
	})
}

func (api *api) changeLang(lang core.Lang) func(m *tb.Message) {
	return func(m *tb.Message) {
		api.execReq(m.Sender, func(ctx context.Context, user *core.User) (message string, err error) {
			if user.CurrentLangState == core.From {
				err = api.app.SetCurrentLangFrom(ctx, *user, lang)
				if err != nil {
					return "", err
				}
				user.CurrentLangFrom = lang
			} else {
				err = api.app.SetCurrentLangTo(ctx, *user, lang)
				if err != nil {
					return "", err
				}
				user.CurrentLangTo = lang
			}
			user.CurrentLangState = core.None
			return fmt.Sprintf("You have chosen %s.", lang.String()), nil
		})
	}
}

func (api *api) changeFromLang() func(m *tb.Message) {
	return func(m *tb.Message) {
		api.execReq(m.Sender, func(ctx context.Context, user *core.User) (message string, err error) {
			err = api.app.SetCurrentLangState(ctx, *user, core.From)
			if err != nil {
				return "", err
			}
			user.CurrentLangState = core.From
			const mess = "Choose a language."

			return mess, nil
		})
	}
}

func (api *api) changeToLang() func(m *tb.Message) {
	return func(m *tb.Message) {
		api.execReq(m.Sender, func(ctx context.Context, user *core.User) (message string, err error) {
			err = api.app.SetCurrentLangState(ctx, *user, core.To)
			if err != nil {
				return "", err
			}
			user.CurrentLangState = core.To
			const mess = "Choose a language."

			return mess, nil
		})
	}
}

type req = func(context.Context, *core.User) (message string, err error)

func (api *api) execReq(sender *tb.User, fn req) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pc, _, _, _ := runtime.Caller(2)
	names := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	method := names[len(names)-1]

	user, err := api.app.CreateOrGetUser(ctx, core.TelegramID(sender.ID), convertLang(sender.LanguageCode))
	if err != nil {
		api.logger.Warn("error:",
			zap.Error(fmt.Errorf("create or get user: %w", err)),
			zap.Int("tgID", sender.ID),
			zap.String("username", sender.Username),
			zap.String("func", method),
		)

	}

	message, err := fn(ctx, user)
	if err != nil {
		api.logger.Warn("error:",
			zap.Error(err),
			zap.Int("tgID", sender.ID),
			zap.String("username", sender.Username),
			zap.String("func", method),
		)

		message = "error"
	}

	b := buttons(*user)
	if user.CurrentLangState == core.From || user.CurrentLangState == core.To {
		b = langButtons
	}

	_, err = api.bot.Send(sender, message, &tb.ReplyMarkup{
		ReplyKeyboard:       b,
		ForceReply:          true,
		ResizeReplyKeyboard: true,
	})
	switch {
	case err != nil:
		api.logger.Warn("error:",
			zap.Error(fmt.Errorf("send message: %w", err)),
			zap.Int("tgID", sender.ID),
			zap.String("username", sender.Username),
			zap.String("func", method),
		)
	default:
		api.logger.Info("success:",
			zap.Int("tgID", sender.ID),
			zap.String("username", sender.Username),
			zap.String("func", method),
		)
	}
}
