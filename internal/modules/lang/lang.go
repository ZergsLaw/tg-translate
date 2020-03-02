package lang

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"github.com/ZergsLaw/tg-translate/internal/core"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type translator struct{}

func (t *translator) Translate(ctx context.Context, text string, from, to core.Lang) (string, error) {
	client, err := translate.NewClient(ctx, option.WithAPIKey("AIzaSyBf_vr_C-Cq1WEjuuGdN46Vk5QcfJksQqE"))
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, parseLang(to), nil)
	if err != nil {
		return "", fmt.Errorf("Translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}

	return resp[0].Text, nil
}

func parseLang(l core.Lang) language.Tag {
	switch l {
	case core.RU:
		return language.Russian
	case core.EN:
		return language.English
	case core.PL:
		return language.Polish
	case core.NL:
		return language.Dutch
	case core.IT:
		return language.Italian
	case core.FR:
		return language.French
	case core.DE:
		return language.German
	case core.PT:
		return language.Portuguese
	case core.ES:
		return language.EuropeanSpanish
	}

	panic(fmt.Sprintf("unknown lang: %d", l))
}

func New() core.Translator {
	return &translator{}
}
