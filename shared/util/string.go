package util

import (
	"MScannot206/shared/def"
	"regexp"
)

var (
	specialCharPatternKorean   = regexp.MustCompile(`[^가-힣a-zA-Z0-9]`)
	specialCharPatternEnglish  = regexp.MustCompile(`[^a-zA-Z0-9]`)
	specialCharPatternJapanese = regexp.MustCompile(`[^\p{Hiragana}\p{Katakana}\p{Han}a-zA-Z0-9]`)
)

func HasSpecialChar(text string, locale def.Locale) bool {
	switch locale {
	case def.LocaleKorean:
		return specialCharPatternKorean.MatchString(text)
	case def.LocaleEnglish:
		return specialCharPatternEnglish.MatchString(text)
	case def.LocaleJapanese:
		return specialCharPatternJapanese.MatchString(text)
	default:
		return specialCharPatternEnglish.MatchString(text)
	}
}
