package telegram

import (
	"text/template"

	"github.com/pkg/errors"
)

type TemplatePack struct {
	rus *template.Template
	eng *template.Template
	tur *template.Template
}

func NewTemplatePack() (*TemplatePack, error) {
	var err error
	var tPack TemplatePack
	dialogTemplateEng := template.New("dialogTemplate").Funcs(template.FuncMap{
		"escape": escapeSimbols,
		"format": formatDateAndMakeLocal,
	})

	dialogTemplateRus := template.New("dialogTemplate").Funcs(template.FuncMap{
		"escape": escapeSimbols,
		"format": formatDateAndMakeLocal,
	})

	dialogTemplateTur := template.New("dialogTemplate").Funcs(template.FuncMap{
		"escape": escapeSimbols,
		"format": formatDateAndMakeLocal,
	})

	tPack.eng, err = dialogTemplateEng.ParseFiles("./templates/dialog_templates_eng.tpl")
	if err != nil {
		return nil, errors.Wrap(err, "Error reading template file")
	}
	tPack.rus, err = dialogTemplateRus.ParseFiles("./templates/dialog_templates_rus.tpl")
	if err != nil {
		return nil, errors.Wrap(err, "Error reading template file")
	}
	tPack.tur, err = dialogTemplateTur.ParseFiles("./templates/dialog_templates_tur.tpl")
	if err != nil {
		return nil, errors.Wrap(err, "Error reading template file")
	}
	return &tPack, err

}
