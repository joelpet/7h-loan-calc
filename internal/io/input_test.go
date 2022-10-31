package io

import (
	"math/big"
	"path"
	"testing"
	"time"
)

func TestReadTransactions(t *testing.T) {
	transactions, err := ReadTransactions(path.Join("..", "testdata", "transaktioner_20220810_20220810.csv"), ';')

	if err != nil {
		t.Errorf("reading transactions: %s", err)
	}

	t.Logf("%+v", transactions)

	if want, got := 1, len(transactions); want != got {
		t.Errorf("want %v, but got %v", want, got)
	}

	if want, got := time.Date(2022, 8, 10, 0, 0, 0, 0, time.UTC), transactions[0].Date; want != got {
		t.Errorf("want %v, but got %v", want, got)
	}
}

func TestParseAmount_Valid(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "0", want: "0.00"},
		{in: "-0", want: "0.00"},
		{in: "1", want: "1.00"},
		{in: "-1", want: "-1.00"},
		{in: "3003", want: "3003.00"},
		{in: "3003,90", want: "3003.90"},
		{in: "3003,904", want: "3003.90"},
		{in: "3003,905", want: "3003.91"},
		{in: "3003,909", want: "3003.91"},
		{in: "3 003,92", want: "3003.92"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseAmount(tt.in)

			if err != nil {
				t.Errorf("parsing amount: %s", err)
			} else if s := got.FloatString(2); s != tt.want {
				t.Errorf("want %+v, but got %+v", tt.want, s)
			}
		})
	}
}

func TestParseAmount_Invalid(t *testing.T) {
	tests := []string{
		"",
		".",
		"3.003,92",
		"3,003,92",
		"3.003.92",
		"3\t003, 92",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			got, err := ParseAmount(tt)

			if got != nil {
				t.Errorf("want no error, but got %+v", got)
			}

			if err == nil {
				t.Error("want error but got nil")
			}
		})
	}
}

func FuzzParseAmount(f *testing.F) {
	tests := []string{"3 003,90", "", " ", ".", ","}

	for _, tt := range tests {
		f.Add(tt)
	}

	f.Fuzz(func(t *testing.T, in string) {
		out, err := ParseAmount(in)

		t.Logf("in=%q, out=%v, err=%v", in, out, err)

		if err != nil && out != nil {
			t.Error("both err and out non-nil")
		}
	})
}

func TestParseAmount_Calculations(t *testing.T) {
	a, err := ParseAmount("1 234,45")
	if err != nil {
		t.Error(err)
	}

	b, err := ParseAmount("678,05")
	if err != nil {
		t.Error(err)
	}

	z := new(big.Rat)
	z.Add(a, b)

	if z.Cmp(big.NewRat(191250, 100)) != 0 {
		t.Error("addition is broken")
	}
}

func TestReadInterestRates(t *testing.T) {
	rates, err := ReadInterestRates(
		path.Join("..", "testdata", "annual_interest_rates.csv"),
		';',
	)
	if err != nil {
		t.Errorf("reading interest rates: %s", err)
	}

	want := []AnnualInterestRate{
		MustNewAnnualInterestRate(2022, time.January, 1, "0.0114"),
		MustNewAnnualInterestRate(2022, time.July, 6, "0.0164"),
		MustNewAnnualInterestRate(2022, time.September, 21, "0.0264"),
		MustNewAnnualInterestRate(2345, time.January, 1, "1.0000"),
	}

	if want, got := len(want), len(rates); want != got {
		t.Errorf("length %d != %d", want, got)
	}

	for i, got := range rates {
		if !want[i].Equal(got) {
			t.Errorf("i=%d, want %s, but got %s", i, want[i], got)
		}
	}
}
