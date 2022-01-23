package translation

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle

func T(id, language string, data interface{}) string {
	localizer := i18n.NewLocalizer(Bundle, language)
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
		TemplateData: data,
	})
}

func Init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	Bundle.MustLoadMessageFile("i18n/en/active.en.toml")
	Bundle.MustLoadMessageFile("i18n/pt/active.pt.toml")
}
