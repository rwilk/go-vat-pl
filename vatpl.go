/*-----------------------------------------------------------------------------
# Name:        GO-VAT-PL
# Purpose:     Verify VAT status in mf.gov.pl database (module)
#
# Author:      Rafal Wilk <rw@pcboot.pl>
#
# Created:     27-09-2020
# Modified:    28-09-2020
# Copyright:   (c) PcBoot 2020
# License:     BSD-new
-----------------------------------------------------------------------------*/

package vatpl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	BLAD = iota
	CZYNNY
	ZWOLNIONY
	NIEZAREJESTROWANY
	NIEZNANY
)

var APIURL = "https://wl-api.mf.gov.pl"

// StatusVAT is enum: BLAD, CZYNNY, ZWOLNIONY, NIEZAREJESTROWANY, NIEZNANY
type StatusVAT int

func (s StatusVAT) String() string {
	switch s {
	case BLAD:
		return "Błąd"
	case CZYNNY:
		return "Czynny"
	case ZWOLNIONY:
		return "Zwolniony"
	case NIEZAREJESTROWANY:
		return "Niezarejestrowany"
	case NIEZNANY:
		return "Nieznany"
	default:
		return fmt.Sprintf("Błąd(val: %d)", s)

	}
}

func (s *StatusVAT) FromString(str string) {
	switch str {
	case "Czynny":
		*s = CZYNNY
	case "Zwolniony":
		*s = ZWOLNIONY
	case "Niezarejestrowany":
		*s = NIEZAREJESTROWANY
	case "":
		*s = NIEZNANY
	default:
		*s = BLAD
	}

}

// VerifyByNIP checks VAT status. Use given date if specified or current if not.
func VerifyByNIP(nip string, date ...interface{}) (status StatusVAT, e error) {

	resource := "/api/search/nip/%s?date=%s"
	tenDigits := regexp.MustCompile(`^\d{10}$`)

	nip = ParseNIP(nip)

	if !tenDigits.MatchString(nip) {
		e = NewVATError("wrong NIP format", true)
		return
	}

	if len(date) > 1 {
		e = NewVATError("wrong number of arguments", true)
		return
	}

	var datestr = time.Now().Format("2006-01-02")

	if len(date) == 1 {
		datestr = dateToStr(date[0])
		if datestr == "" {
			e = NewVATError("wrong date format: use YYYY-MM-dd or time.Time", true)
			return
		}
	}

	URL := APIURL + fmt.Sprintf(resource, nip, datestr)

	client := &http.Client{}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		e = NewVATError(err.Error(), false)
		return
	}

	// acting like web browser to bypass limits
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36")
	req.Header.Set("Accept-Language", "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	//req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		e = NewVATError(err.Error(), false)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e = NewVATError(err.Error(), false)
		return
	}

	var vatresponse VATResponse

	if err = json.Unmarshal(body, &vatresponse); err != nil {
		e = NewVATError(err.Error(), false)
		return
	}

	if vatresponse.Code != "" {
		switch vatresponse.Code {
		case "WL-100", "WL-191", "WL-195", "WL-196":
			e = NewVATError(fmt.Sprintf("%s: %s", vatresponse.Code, vatresponse.Message), false)
		default:
			e = NewVATError(fmt.Sprintf("%s: %s", vatresponse.Code, vatresponse.Message), true)
		}
		return
	}

	status.FromString(vatresponse.Result.Subject.StatusVat)

	return
}

// ParseNIP - trim, remove "-" and return string
func ParseNIP(nip string) string {
	nip = strings.Trim(nip, " \t\n\r")
	return strings.ReplaceAll(nip, "-", "")
}

// dateToStr try to convert interface-type date to string YYYY-MM-dd
// Returns empty string on error
func dateToStr(i interface{}) string {
	switch v := i.(type) {
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return ""
		}
		return t.Format("2006-01-02")
	case time.Time:
		return v.Format("2006-01-02")
	default:
		return ""
	}
}
