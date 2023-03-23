package telegram

import (
	"text/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type LanguagePack struct {
	HoursKeyboard       tgbotapi.ReplyKeyboardMarkup
	MenuKeyboard        tgbotapi.ReplyKeyboardMarkup
	ConfirmKeyboard     tgbotapi.ReplyKeyboardMarkup
	SettingsKeyboard    tgbotapi.ReplyKeyboardMarkup
	LanguageKeyboard    tgbotapi.ReplyKeyboardMarkup
	CityKeyboard        tgbotapi.ReplyKeyboardMarkup
	FethiyeKeyboard     tgbotapi.ReplyKeyboardMarkup
	DalamanKeyboard     tgbotapi.ReplyKeyboardMarkup
	KavaklidereKeyboard tgbotapi.ReplyKeyboardMarkup
	KoycegizKeyboard    tgbotapi.ReplyKeyboardMarkup
	MarmarisKeyboard    tgbotapi.ReplyKeyboardMarkup
	MenteseKeyboard     tgbotapi.ReplyKeyboardMarkup
	MilasKeyboard       tgbotapi.ReplyKeyboardMarkup
	SeydikemerKeyboard  tgbotapi.ReplyKeyboardMarkup
	UlaKeyboard         tgbotapi.ReplyKeyboardMarkup
	YataganKeyboard     tgbotapi.ReplyKeyboardMarkup
	OrtacaKeyboard      tgbotapi.ReplyKeyboardMarkup
	DatcaKeyboard       tgbotapi.ReplyKeyboardMarkup
	BodrumKeyboard      tgbotapi.ReplyKeyboardMarkup
	Template            *template.Template
}

type Languages struct {
	Eng LanguagePack
	Rus LanguagePack
	Tur LanguagePack
}

func NewLanguages() (Languages, error) {
	//parsing the template file
	t, err := NewTemplatePack()
	if err != nil {
		return Languages{}, err
	}
	English := LanguagePack{
		HoursKeyboard:       HoursKeyboardEng,
		MenuKeyboard:        MenuKeyboardEng,
		ConfirmKeyboard:     ConfirmKeyboardEng,
		SettingsKeyboard:    SettingsKeyboardEng,
		LanguageKeyboard:    LanguageKeyboardEng,
		CityKeyboard:        CityKeyboardEng,
		FethiyeKeyboard:     FethiyeKeyboardEng,
		DalamanKeyboard:     DalamanKeyboardEng,
		KavaklidereKeyboard: KavaklidereKeyboardEng,
		KoycegizKeyboard:    KoycegizKeyboardEng,
		MarmarisKeyboard:    MarmarisKeyboardEng,
		MenteseKeyboard:     MenteseKeyboardEng,
		MilasKeyboard:       MilasKeyboardEng,
		SeydikemerKeyboard:  SeydikemerKeyboardEng,
		UlaKeyboard:         UlaKeyboardEng,
		YataganKeyboard:     YataganKeyboardEng,
		OrtacaKeyboard:      OrtacaKeyboardEng,
		DatcaKeyboard:       DatcaKeyboardEng,
		BodrumKeyboard:      BodrumKeyboardEng,
		Template:            t.eng,
	}

	Russian := LanguagePack{
		HoursKeyboard:       HoursKeyboardRus,
		MenuKeyboard:        MenuKeyboardRus,
		ConfirmKeyboard:     ConfirmKeyboardRus,
		SettingsKeyboard:    SettingsKeyboardRus,
		LanguageKeyboard:    LanguageKeyboardRus,
		CityKeyboard:        CityKeyboardRus,
		FethiyeKeyboard:     FethiyeKeyboardRus,
		DalamanKeyboard:     DalamanKeyboardRus,
		KavaklidereKeyboard: KavaklidereKeyboardRus,
		KoycegizKeyboard:    KoycegizKeyboardRus,
		MarmarisKeyboard:    MarmarisKeyboardRus,
		MenteseKeyboard:     MenteseKeyboardRus,
		MilasKeyboard:       MilasKeyboardRus,
		SeydikemerKeyboard:  SeydikemerKeyboardRus,
		UlaKeyboard:         UlaKeyboardRus,
		YataganKeyboard:     YataganKeyboardRus,
		OrtacaKeyboard:      OrtacaKeyboardRus,
		DatcaKeyboard:       DatcaKeyboardRus,
		BodrumKeyboard:      BodrumKeyboardRus,
		Template:            t.rus,
	}

	Turkish := LanguagePack{
		HoursKeyboard:       HoursKeyboardTur,
		MenuKeyboard:        MenuKeyboardTur,
		ConfirmKeyboard:     ConfirmKeyboardTur,
		SettingsKeyboard:    SettingsKeyboardTur,
		LanguageKeyboard:    LanguageKeyboardTur,
		CityKeyboard:        CityKeyboardTur,
		FethiyeKeyboard:     FethiyeKeyboardTur,
		DalamanKeyboard:     DalamanKeyboardTur,
		KavaklidereKeyboard: KavaklidereKeyboardTur,
		KoycegizKeyboard:    KoycegizKeyboardTur,
		MarmarisKeyboard:    MarmarisKeyboardTur,
		MenteseKeyboard:     MenteseKeyboardTur,
		MilasKeyboard:       MilasKeyboardTur,
		SeydikemerKeyboard:  SeydikemerKeyboardTur,
		UlaKeyboard:         UlaKeyboardTur,
		YataganKeyboard:     YataganKeyboardTur,
		OrtacaKeyboard:      OrtacaKeyboardTur,
		DatcaKeyboard:       DatcaKeyboardTur,
		BodrumKeyboard:      BodrumKeyboardTur,
		Template:            t.tur,
	}

	return Languages{
		Eng: English,
		Rus: Russian,
		Tur: Turkish,
	}, err

}
