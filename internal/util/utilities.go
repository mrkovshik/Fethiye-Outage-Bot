package util

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Normalize(s string) ([]string, error) {
	tr := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(tr, s)
	if err != nil {
		err = errors.Wrap(err, "Error transforming input query")
		return []string{}, err
	}
	re, err := regexp.Compile(`[\p{L}\d_]+`)
	if err != nil {
		err = errors.Wrap(err, "Error compiling regexp input query")
		return []string{}, err
	}
	res := re.FindAllString(strings.ToLower(output), -1)
	return res, err
}
