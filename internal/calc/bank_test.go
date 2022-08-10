package calc

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"gitlab.joelpet.se/joelpet/7h-loan-calc/internal/io"
)

func TestProcess(t *testing.T) {
	bank := NewBank([]io.Transaction{
		io.MustNewTransaction(2022, 1, 2, "3003.90"),
		io.MustNewTransaction(2022, 2, 3, "4000.50"),
		io.MustNewTransaction(2022, 3, 1, "2000.70"),
	}, []io.AnnualInterestRate{
		io.MustNewAnnualInterestRate(2018, 1, 1, "0.0099"),
		io.MustNewAnnualInterestRate(2021, 6, 27, "0.0114"),
		io.MustNewAnnualInterestRate(2022, 7, 31, "0.0164"),
		io.MustNewAnnualInterestRate(2022, 12, 31, "0.0179"),
	})

	loan := Loan{
		balance:  mustBigRatFromString("100000"),
		interest: mustBigRatFromString("0"),
	}
	day := time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)

	want := Loan{
		balance:  mustBigRatFromString("100000"),
		interest: mustBigRatFromString("2.75"),
	}

	got := bank.Process(day, loan)

	if cmp := want.balance.Cmp(got.balance); cmp != 0 {
		t.Errorf("want balance %s, but got %s (cmp=%d)", want.balance, got.balance, cmp)
	}

	if cmp := want.interest.Cmp(got.interest); cmp != 0 {
		t.Errorf("want interest %s, but got %s (cmp=%d)", want.interest, got.interest, cmp)
	}
}

func TestTransactionsAmount(t *testing.T) {
	bank := NewBank([]io.Transaction{
		io.MustNewTransaction(2022, 1, 2, "3003.90"),
		io.MustNewTransaction(2022, 2, 3, "4000.50"),
		io.MustNewTransaction(2022, 3, 1, "2000.70"),
	}, []io.AnnualInterestRate{})

	locNY, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Errorf("loading location: %s", err)
	}

	tests := []struct {
		name       string
		date       time.Time
		wantAmount *big.Rat
	}{
		{
			name:       "no transaction",
			date:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			wantAmount: mustBigRatFromString("0"),
		},
		{
			name:       "transaction midnight",
			date:       time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			wantAmount: mustBigRatFromString("3003.90"),
		},
		{
			name:       "transaction noon",
			date:       time.Date(2022, 1, 2, 12, 0, 0, 0, time.UTC),
			wantAmount: mustBigRatFromString("3003.90"),
		},
		{
			name:       "transaction New York timezone",
			date:       time.Date(2022, 1, 2, 0, 0, 0, 0, locNY),
			wantAmount: mustBigRatFromString("3003.90"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := bank.transactionsAmount(tt.date)

			if a.Cmp(tt.wantAmount) != 0 {
				t.Errorf("want amount %q, got %s", tt.wantAmount, a)
			}
		})
	}
}

func TestAnnualInterestRate(t *testing.T) {
	bank := NewBank([]io.Transaction{}, []io.AnnualInterestRate{
		io.MustNewAnnualInterestRate(2018, 1, 1, "0.0099"),
		io.MustNewAnnualInterestRate(2021, 6, 27, "0.0114"),
		io.MustNewAnnualInterestRate(2022, 7, 31, "0.0164"),
		io.MustNewAnnualInterestRate(2022, 12, 31, "0.0179"),
	})

	tests := []struct {
		name     string
		date     time.Time
		wantOK   bool
		wantRate *big.Rat
	}{
		{
			name:   "no rate",
			date:   time.Date(1999, 5, 19, 0, 0, 0, 0, time.UTC),
			wantOK: false,
		},
		{
			name:     "oldest rate on day of change",
			date:     time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			wantOK:   true,
			wantRate: mustBigRatFromString("0.0099"),
		},
		{
			name:     "oldest rate between days of change",
			date:     time.Date(2019, 2, 28, 0, 0, 0, 0, time.UTC),
			wantOK:   true,
			wantRate: mustBigRatFromString("0.0099"),
		},
		{
			name:     "mid-rate on day of change",
			date:     time.Date(2021, 6, 27, 0, 0, 0, 0, time.UTC),
			wantOK:   true,
			wantRate: mustBigRatFromString("0.0114"),
		},
		{
			name:     "mid-rate between days of change",
			date:     time.Date(2021, 10, 15, 0, 0, 0, 0, time.UTC),
			wantOK:   true,
			wantRate: mustBigRatFromString("0.0114"),
		},
		{
			name:     "most recent rate on day of change",
			date:     time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
			wantOK:   true,
			wantRate: mustBigRatFromString("0.0179"),
		},
		{
			name:     "most recent rate after day of change",
			date:     time.Date(2023, 2, 19, 0, 0, 0, 0, time.UTC),
			wantOK:   true,
			wantRate: mustBigRatFromString("0.0179"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate, ok := bank.annualInterestRate(tt.date)

			if tt.wantOK != ok {
				t.Errorf("want OK %t, but got %t", tt.wantOK, ok)
			}

			if ok && tt.wantRate.Cmp(rate) != 0 {
				t.Errorf("want rate %s, but got %s", tt.wantRate, rate)
			}
		})
	}
}

// TestAnnualToDailyAppliedToBalance tests annualToDaily by
// calculating the annual interest cost backwards starting with the
// daily interest rate and a known balance, scaling it up first to
// a monthly cost, and finally to an annual cost. The annual cost is
// easily calculated by hand as (yearly) interest rate times the
// balance.
func TestAnnualToDailyAppliedToBalance(t *testing.T) {
	balance := mustBigRatFromString("100")
	rate := mustBigRatFromString("0.02")
	daysInMonth := 30
	daysInM := new(big.Rat).SetInt64(int64(daysInMonth))
	monthsInY := new(big.Rat).SetInt64(12)

	dailyRate := annualToDaily(rate, daysInMonth)

	dailyCost := new(big.Rat).Mul(balance, dailyRate)
	monthlyCost := new(big.Rat).Mul(dailyCost, daysInM)
	annualCost := new(big.Rat).Mul(monthlyCost, monthsInY)

	wantAnnualCost := new(big.Rat).Mul(balance, rate)

	if cmp := wantAnnualCost.Cmp(annualCost); cmp != 0 {
		t.Errorf("want %s, but got %s",
			wantAnnualCost.String(),
			annualCost.String())
	}
}

func mustBigRatFromString(s string) *big.Rat {
	if f, ok := new(big.Rat).SetString(s); ok {
		return f
	} else {
		panic(fmt.Errorf("set big rat to %q", s))
	}
}
