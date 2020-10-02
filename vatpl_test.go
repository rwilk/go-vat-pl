package vatpl

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// KGHM Polska Miedź
const CzynnyNIP = "692-00-00-013"

// Computer generated
const NieznanyNIP = "375-17-84-446"

func TestVerifyByNIP_DirtyFormat(t *testing.T) {
	// dirty format
	nip := " " + CzynnyNIP + "\t"
	want := StatusVAT(CZYNNY)
	got, err := VerifyByNIP(nip)
	if err != nil {
		t.Errorf(err.Error())
	}
	if got != want {
		t.Errorf("VerifyByNIP = %q, want %q", got, want)
	}
}

func TestVerifyByNIP_DateAsString(t *testing.T) {
	nip := CzynnyNIP
	want := StatusVAT(CZYNNY)

	got, err := VerifyByNIP(nip, time.Now().Format("2006-01-02"))

	if err != nil {
		t.Errorf(err.Error())
	}
	if got != want {
		t.Errorf("VerifyByNIP = %q, want %q", got, want)
	}
}

func TestVerifyByNIP_DateAsTime(t *testing.T) {
	nip := CzynnyNIP
	want := StatusVAT(CZYNNY)

	got, err := VerifyByNIP(nip, time.Now())

	if err != nil {
		t.Errorf(err.Error())
	}
	if got != want {
		t.Errorf("VerifyByNIP = %q, want %q", got, want)
	}
}

func TestVerifyByNIP_DateFromFuture(t *testing.T) {
	nip := CzynnyNIP
	want := StatusVAT(BLAD)

	got, _ := VerifyByNIP(nip, time.Now().Add(time.Hour*24))

	if got != want {
		t.Errorf("VerifyByNIP = %q, want %q", got, want)
	}
}

func TestVerifyByNIP_DateFromPast(t *testing.T) {
	nip := CzynnyNIP
	want := StatusVAT(CZYNNY)

	got, _ := VerifyByNIP(nip, time.Now().Add(time.Hour*72*-1))

	if got != want {
		t.Errorf("VerifyByNIP = %q, want %q", got, want)
	}
}

func TestVerifyByNIP_DateWrongString(t *testing.T) {
	nip := CzynnyNIP
	want := StatusVAT(BLAD)

	got, err := VerifyByNIP(nip, time.Now().Format("02-01-2006"))

	if err == nil {
		t.Errorf("VerifyByNIP error nil")
	}

	if got != want {
		t.Errorf("VerifyByNIP = %q, want %q", got, want)
	}
}

func TestVerifyByNIP_WrongNIPFormat(t *testing.T) {
	nip := "999" + CzynnyNIP[3:]
	want := StatusVAT(BLAD)

	got, _ := VerifyByNIP(nip)

	if got != want {
		t.Errorf("VerifyByNIP = %q, want %q", got, want)
	}
}

func TestVerifyByNIP_UnknownNIP(t *testing.T) {
	nip := NieznanyNIP
	want := StatusVAT(NIEZNANY)

	got, _ := VerifyByNIP(nip)

	if got != want {
		t.Errorf("VerifyByNIP = %q, want %q", got, want)
	}
}

func TestVATError_IsPermanent(t *testing.T) {
	// test 1 - expected permanent == true
	nip := "1234567891"
	want := true

	_, err := VerifyByNIP(nip)

	if err == nil {
		t.Errorf("VATError (test1) expected, got nil")
		return
	}

	e, ok := err.(*VATError)

	if !ok {
		t.Errorf("VATError (test1) asseration failed")
		return
	}

	got := e.IsPermanent()

	if got != want {
		t.Errorf("VATError (test1) = %v, want %v", got, want)
	}

	// test 2: expected permanent == false
	var temp_apiurl = APIURL
	APIURL = APIURL + ".wrong.address"

	defer func() {
		APIURL = temp_apiurl
	}()

	want = false
	_, err = VerifyByNIP(CzynnyNIP)

	if err == nil {
		t.Errorf("VATError (test2) expected, got nil")
		return
	}

	e, ok = err.(*VATError)

	if !ok {
		t.Errorf("VATError (test2) asseration failed")
		return
	}

	got = e.IsPermanent()

	if got != want {
		t.Errorf("VATError (test2) = %v, want %v", got, want)
	}

}

func TestVerifyByNIPRetry(t *testing.T) {
	nip := CzynnyNIP
	want := StatusVAT(CZYNNY)
	got, err := VerifyByNIPRetry(nip)
	if err != nil {
		t.Errorf(err.Error())
	}
	if got != want {
		t.Errorf("VerifyByNIPRetry = %q, want %q", got, want)
	}
}

func TestNIPPortions(t *testing.T) {
	nips := make(map[string]string, 65)

	// generate pseudo NIPs (65 pcs)
	for i := 1; i <= 65; i++ {
		s := fmt.Sprintf("%010d", i)
		nips[s] = s
	}

	portions := NIPPortions(nips)

	if len(portions) != 3 {
		t.Errorf("NIPPortions = %d, want %d (slice length)", len(portions), 3)
		return
	}

	p1 := strings.Split(portions[0], ",")
	if len(p1) != 30 {
		t.Errorf("NIPPortions = %d, want %d (first portion)", len(p1), 30)
		return
	}

	p2 := strings.Split(portions[1], ",")
	if len(p2) != 30 {
		t.Errorf("NIPPortions = %d, want %d (second portion)", len(p2), 30)
		return
	}

	p3 := strings.Split(portions[2], ",")
	if len(p3) != 5 {
		t.Errorf("NIPPortions = %d, want %d (third portion)", len(p3), 5)
		return
	}
}

func TestSplitNIPS(t *testing.T) {

	nips := []string{CzynnyNIP, "1122334455", "12345", NieznanyNIP, "1234567890"}
	nMap := make(map[string]string)

	for _, n := range nips {
		nMap[n] = ParseNIP(n)
	}

	good, bad := SplitNIPS(nMap)

	if len(good) != 2 {
		t.Errorf("SplitNIPS = %d, want %d (good map length)", len(good), 2)
		return
	}

	if len(bad) != 3 {
		t.Errorf("SplitNIPS = %d, want %d (bad map length)", len(bad), 3)
		return
	}
}

func TestVerifyByNIPBulkRetry_ContainsWrongNIP(t *testing.T) {
	nips := []string{CzynnyNIP, "1122334455", "12345", NieznanyNIP, "1234567890"}

	results, err := VerifyByNIPBulkRetry(nips)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if len(results) != len(nips) {
		t.Errorf("VerifyByNIPBulRetry = %d, want %d (nips in result)", len(results), len(nips))
		return
	}
}
