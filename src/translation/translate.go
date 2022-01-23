package translation

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle

func T(id, language string, data ...interface{}) string {
	localizer := i18n.NewLocalizer(Bundle, language)
	config := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
	}
	if len(data) > 0 {
		config.TemplateData = data[0]
	}
	return localizer.MustLocalize(config)
}

func init() {
	Bundle = i18n.NewBundle(language.Portuguese)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, lang := range [...]string{"pt"} {
		for _, file := range [...]string{"bot", "commands"} {
			Bundle.MustLoadMessageFile(fmt.Sprintf("i18n/%s/%s.%s.toml", lang, file, lang))
		}
	}
}
