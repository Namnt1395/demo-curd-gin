package i18n

import (
	"bytes"
	"demo-curd/config"
	"demo-curd/util/constant"
	"encoding/json"
	"golang.org/x/text/language"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"html/template"
	"io/ioutil"
	"path"
	"strings"
)

type I18n struct {
	Bundle       *i18n.Bundle
	MapLocalizer map[string]*i18n.Localizer
}

func NewI18n(c config.Config) (*I18n, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	mapLocalizer := make(map[string]*i18n.Localizer)
	files, err := ioutil.ReadDir("i18n")
	if err != nil {
		return nil, err
	}
	for _, lang := range c.I18n.Langs {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), lang+".json") && !file.IsDir() {
				bundle.MustLoadMessageFile(path.Join("i18n", file.Name()))
			}
		}
		mapLocalizer[lang] = i18n.NewLocalizer(bundle, lang)
	}

	return &I18n{
		Bundle:       bundle,
		MapLocalizer: mapLocalizer,
	}, nil
}

func (r *I18n) MustLocalize(lang string, msgId string, templateData map[string]string, defaultMsg ...string) string {
	var localize *i18n.Localizer
	localize, ok := r.MapLocalizer[lang]
	if !ok {
		localize = r.MapLocalizer[constant.DefaultLang]
	}
	ret, err := localize.Localize(&i18n.LocalizeConfig{
		MessageID:    msgId,
		TemplateData: templateData,
	})
	if err != nil {
		if len(defaultMsg) > 0 {
			return bindingTemplate(msgId, defaultMsg[0], templateData)
		} else {
			return msgId
		}
	}
	return ret
}

func bindingTemplate(name string, msg string, data interface{}) string {
	tmpl, err := template.New(name).Parse(msg)
	if err != nil {
		return msg
	}
	var buffer = bytes.NewBuffer(nil)
	if err = tmpl.Execute(buffer, data); err != nil {
		return msg
	}
	return buffer.String()
}
