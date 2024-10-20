// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package function

import (
	"fmt"
	"github.com/jrnd-io/jrv2/pkg/state"
	"strconv"
	"strings"
	"text/template"

	"github.com/jrnd-io/jrv2/pkg/random"
	"github.com/squeeze69/generacodicefiscale"

	"github.com/rs/zerolog/log"
)

func init() {
	AddFuncs(template.FuncMap{
		// people related utilities
		"cf":             CodiceFiscale,
		"company":        Company,
		"email":          Email,
		"email_provider": EmailProvider,
		"email_work":     WorkEmail,
		"gender":         Gender,
		"middlename":     Middlename,
		"name":           Name,
		"name_m":         NameM,
		"name_f":         NameF,
		"ssn":            Ssn,
		"surname":        Surname,
		"user":           User,
		"username":       Username,
	})
}

const (
	CompanyMap      = "company"
	NameFMap        = "nameF"
	NameMMap        = "nameM"
	MailProviderMap = "mail_provider"
	SurnameMap      = "surname"
)

// CodiceFiscale return a valid Italian Codice Fiscale
func CodiceFiscale() string {

	name := state.GetSharedState().Value("_name").(string)
	surname := state.GetSharedState().Value("_surname").(string)
	gender := state.GetSharedState().Value("_gender").(string)
	birthdate := state.GetSharedState().Value("_birthdate").(string)
	city := state.GetSharedState().Value("_city").(string)

	if name == "" {
		name = Name()
	}
	if surname == "" {
		surname = Surname()
	}
	if gender == "" {
		gender = Gender()
	}
	if birthdate == "" {
		birthdate = BirthDate(18, 75)
	}
	if city == "" {
		city = City()
	}

	if city == "Bolzano" {
		city = "Bolzano/Bozen"
	}
	if city == "Reggio Emilia" {
		city = "Reggio nell'Emilia"
	}
	if city == "Reggio Calabria" {
		city = "Reggio di Calabria"
	}

	codicecitta, erc := generacodicefiscale.CercaComune(city)
	if erc != nil {
		log.Fatal().Err(erc).Msg("Error in searching city")
	}
	cf, erg := generacodicefiscale.Genera(surname, name, gender, codicecitta.Codice, birthdate)
	if erg != nil {
		log.Fatal().Err(erg).Msg("Error in generating Codice Fiscale")
	}
	return cf
}

// Company returns a random Company Name
func Company() string {
	c := Word(CompanyMap)
	state.GetSharedState().Ctx.Store("_company", c)
	return c
}

// WorkEmail returns a random work email.
func WorkEmail() string {
	name := state.GetSharedState().Value("_name").(string)
	surname := state.GetSharedState().Value("_surname").(string)
	company := state.GetSharedState().Value("_company").(string)

	if name == "" {
		name = Name()
	}
	if surname == "" {
		surname = Surname()
	}
	if company == "" {
		company = Company()
	}
	company = strings.ReplaceAll(company, " ", "")
	return fmt.Sprintf("%s.%s@%s.com", strings.ToLower(name), strings.ToLower(surname), strings.ToLower(company))
}

// Email returns a random email.
func Email() string {
	name := state.GetSharedState().Value("_name").(string)
	surname := state.GetSharedState().Value("_surname").(string)
	provider := Word(MailProviderMap)

	if name == "" {
		name = Name()
	}
	if surname == "" {
		surname = Surname()
	}

	return fmt.Sprintf("%s.%s@%s", strings.ToLower(name), strings.ToLower(surname), strings.ToLower(provider))
}

// EmailProvider returns a random email provider
func EmailProvider() string {
	return Word(MailProviderMap)
}

// Gender returns a random gender. Note: it gets the gender context automatically setup by previous name calls
func Gender() string {
	g := state.GetSharedState().Value("_gender").(string)
	if g == "" {
		gender := []string{"M", "F"}
		g = gender[random.Random.IntN(len(gender))]
		state.GetSharedState().Ctx.Store("_gender", g)
	}
	return g
}

// Middlename returns a random Middlename
func Middlename() string {
	middles := []string{"M", "J", "K", "P", "T", "S"}
	return middles[random.Random.IntN(len(middles))]
}

// Name returns a random Name (male/female)
func Name() string {
	s := random.Random.IntN(2)
	if s == 0 {
		return NameM()
	}

	return NameF()
}

// NameM returns a random male Name
func NameM() string {
	name := Word(NameMMap)
	state.GetSharedState().Ctx.Store("_name", name)
	state.GetSharedState().Ctx.Store("_gender", "M")
	return name
}

// NameF returns a random female Name
func NameF() string {
	name := Word(NameFMap)
	state.GetSharedState().Ctx.Store("_name", name)
	state.GetSharedState().Ctx.Store("_gender", "F")
	return name
}

// Ssn return a valid Social Security Number id
func Ssn() string {
	first := random.Random.IntN(899) + 1
	second := random.Random.IntN(99) + 1
	third := random.Random.IntN(9999) + 1
	return fmt.Sprintf("%03d-%02d-%04d", first, second, third)
}

// Surname returns a random Surname
func Surname() string {
	s := Word(SurnameMap)
	state.GetSharedState().Ctx.Store("_surname", "F")
	return s
}

// Username returns a random Username using Name, Surname
func Username(firstName string, lastName string) string {

	firstName = strings.ToLower(firstName)
	lastName = strings.ToLower(lastName)

	separators := []string{".", "-", "", "_", "."}
	separator := separators[random.Random.IntN(len(separators))]
	onlyInitialForName := (random.Random.IntN(2)) != 0
	onlyInitialForSurname := (random.Random.IntN(2)) != 0
	useSurname := (random.Random.IntN(2)) != 0

	if onlyInitialForName {
		firstName = firstName[:1]
		onlyInitialForSurname = false
		useSurname = true
	}
	if onlyInitialForSurname {
		lastName = lastName[:1]
	}
	if !(useSurname) {
		lastName = ""
		separator = ""
	}
	return fmt.Sprintf("%s%s%s", firstName, separator, lastName)

}

// User returns a random Username using Name, Surname and a length
func User(firstName string, lastName string, size int) string {

	var name string

	useSurname := (random.Random.IntN(2)) != 0
	shuffleName := (random.Random.IntN(2)) != 0
	if useSurname || len(name) < size {
		name = firstName + lastName
	} else {
		name = firstName
	}

	if shuffleName {
		nameRunes := []rune(name)
		random.Random.Shuffle(len(nameRunes), func(i, j int) {
			nameRunes[i], nameRunes[j] = nameRunes[j], nameRunes[i]
		})
		name = string(nameRunes)
	}

	var username string
	for i := 0; i < size && i < len(name); i++ {
		username += string(name[i])
	}

	username += strconv.Itoa(50 + random.Random.IntN(49))

	return username
}
