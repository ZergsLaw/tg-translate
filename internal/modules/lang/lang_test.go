package lang

import (
	"context"
	"github.com/ZergsLaw/tg-translate/internal/core"
	"log"
	"testing"
)

func TestTranslator_Translate(t *testing.T) {
	t.Parallel()

	translator := &translator{}

	text, err := translator.Translate(context.Background(), "Hello world!", core.EN, core.RU)
	log.Println(text, err)
}
