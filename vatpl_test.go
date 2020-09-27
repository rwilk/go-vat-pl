package vatpl

import (
	"testing"
	"time"
)

// KGHM Polska Miedź
const CzynnyNIP = "692-00-00-013"

// Computer generated
const NieznanyNIP = "375-17-84-446"

func TestVerifyByNIP_DirtyFormat(t *testing.T) {
	// KGHM Polska Miedź
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
	// KGHM Polska Miedź
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
