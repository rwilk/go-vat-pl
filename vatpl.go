/*-----------------------------------------------------------------------------
# Name:        GO-VAT-PL
# Purpose:     Verify VAT status in mf.gov.pl database (module)
#
# Author:      Rafal Wilk <rw@pcboot.pl>
#
# Created:     27-09-2020
# Modified:    07-10-2020
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
	"strconv"
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
var RETRYCOUNT = 5

var tenDigits = regexp.MustCompile(`^\d{10}$`)

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

// VerifyByNIPRetry same as VeryfiByNIP but retry on non-permanent errors
func VerifyByNIPRetry(nip string, date ...interface{}) (status StatusVAT, e error) {
	var sleepSec = 1

	for rc := 0; rc < RETRYCOUNT; rc++ {
		status, e = VerifyByNIP(nip, date...)

		if e == nil {
			return
		}

		ev, ok := e.(*VATError)
		if !ok {
			return
		}

		if ev.Permanent {
			return
		}

		// skip last sleep
		if rc < RETRYCOUNT-1 {
			time.Sleep(time.Duration(sleepSec) * time.Second)
			sleepSec = sleepSec * 2
		}
	}

	return
}

// VerifyByNIP checks VAT status. Use given date if specified or current if not.
func VerifyByNIP(nip string, date ...interface{}) (status StatusVAT, e error) {

	resource := "/api/search/nip/%s?date=%s"

	nip = ParseNIP(nip)

	if !IsValidNIP(nip) {
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

// VerifyByNIPBulkRetry same as VeryfiByNIPBulk but retry on non-permanent errors
func VerifyByNIPBulkRetry(nips []string, date ...interface{}) (statuses map[string]StatusVAT, e error) {
	var sleepSec = 1

	for rc := 0; rc < RETRYCOUNT; rc++ {
		statuses, e = VerifyByNIPBulk(nips, date...)

		if e == nil {
			return
		}

		ev, ok := e.(*VATError)
		if !ok {
			return
		}

		if ev.Permanent {
			return
		}

		// skip last sleep
		if rc < RETRYCOUNT-1 {
			time.Sleep(time.Duration(sleepSec) * time.Second)
			sleepSec = sleepSec * 2
		}

	}

	return
}

// VerifyByNIPBulk checks VAT status for slice of NIPs. Use given date if specified or current if not.
func VerifyByNIPBulk(nips []string, date ...interface{}) (statuses map[string]StatusVAT, e error) {
	statuses = make(map[string]StatusVAT, len(nips))
	resource := "/api/search/nips/%s?date=%s"

	// nips in "clean" format
	nMap := make(map[string]string)

	for _, n := range nips {
		nMap[n] = ParseNIP(n)
	}
	//nMapInv := invertMap(nMap)

	nMapGood, nMapBad := SplitNIPS(nMap)

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

	nips = NIPPortions(nMapGood)

	// sets all "good" nips as NIEZNANY in case they are missing in API response
	for k, _ := range nMapGood {
		statuses[k] = NIEZNANY
	}

	// add bad nips
	for k, _ := range nMapBad {
		statuses[k] = BLAD
	}

	for _, portion := range nips {
		URL := APIURL + fmt.Sprintf(resource, portion, datestr)

		client := &http.Client{}

		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			e = NewVATError(err.Error(), false)
			return
		}

		// acting like web browser to bypass limits
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36")
		req.Header.Set("Accept-Language", "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7")
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

		for _, s := range vatresponse.Result.Subjects {
			var status StatusVAT
			status.FromString(s.StatusVat)
			//statuses[nMapInv[s.Nip]] = status
			statuses = setForAllKeys(statuses, s.Nip, status)
		}
	}
	return
}

// ParseNIP - trim, remove "-" and return string
func ParseNIP(nip string) string {
	nip = strings.Trim(nip, " \t\n\r")
	return strings.ReplaceAll(nip, "-", "")
}

// NIPPortions converts parsed NIPs to 30 elements comma-separated strings
func NIPPortions(nMap map[string]string) []string {
	var (
		portions []string
		s        string
		i        int
		first    bool = true
	)

	for _, v := range nMap {
		if first {
			s = v
			first = false
		} else {
			s += "," + v
		}
		i++
		if i == 30 {
			portions = append(portions, s)
			first = true
			i = 0
		}
	}
	// unfinished portion
	portions = append(portions, s)

	return portions
}

// IsValidNIP checks if NIP is valid
func IsValidNIP(nip string) bool {

	if !tenDigits.MatchString(nip) {
		return false
	}

	var weights = []int{6, 5, 7, 2, 3, 4, 5, 6, 7}
	var sum int

	for i, s := range nip {
		if i == 9 {
			break
		}
		n, _ := strconv.Atoi(string(s))
		sum += n * weights[i]
	}

	crc, _ := strconv.Atoi(string(nip[9]))

	if crc == sum%11 {
		return true
	}

	return false
}

// SplitNIPS to two maps contains valid (good) and invalid (bad) NIPs
func SplitNIPS(nMap map[string]string) (good, bad map[string]string) {
	good = make(map[string]string)
	bad = make(map[string]string)

	for k, v := range nMap {
		if IsValidNIP(v) {
			good[k] = v
		} else {
			bad[k] = v
		}
	}

	return
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

// invertMap swaps values with keys
func invertMap(m map[string]string) map[string]string {
	inverted := make(map[string]string, len(m))

	for k, v := range m {
		inverted[v] = k
	}

	return inverted
}

// setForAllKeys - sets status for all keys with given nip (all formats)
func setForAllKeys(nMap map[string]StatusVAT, parsedNIP string, status StatusVAT) map[string]StatusVAT {
	for k, _ := range nMap {
		if pk := ParseNIP(k); pk == parsedNIP {
			nMap[k] = status
		}
	}

	return nMap
}
