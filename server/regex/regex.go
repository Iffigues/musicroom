package regex

import (
	"regexp"
	"unicode"
)

type Validateur struct{
	Name *regexp.Regexp
	Email *regexp.Regexp
	Username *regexp.Regexp
}

var (
	Verifie Validateur
)

func (v Validateur)IsEmail(d string) (boolean bool){
	return v.Email.MatchString(d) && d != ""
}

func (v Validateur)IsName(d string) (boolean bool){
	return v.Name.MatchString(d) && d != ""
}

func (v Validateur) IsUsername(d string)(boolean bool) {
	return v.Username.MatchString(d) && d != ""
}

func (v Validateur)ValidPassword(s string) bool {
	if len(s) < 8 {
		return false
	}
	next:
	for _, classes := range map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"lower case": {unicode.Lower},
		"numeric":    {unicode.Number, unicode.Digit},
		"special":    {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
	} {
		for _, r := range s {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return false
	}
	return true
}

func init () {
	Verifie.Name = regexp.MustCompile("^[a-z ,.'-]+")
	Verifie.Email = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	Verifie.Username = regexp.MustCompile("^[a-zA-Z0-9_.]*$")
}
