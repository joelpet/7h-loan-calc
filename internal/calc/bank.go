package calc

import (
	"math/big"
	"sort"
	"time"

	"gitlab.joelpet.se/joelpet/7h-loan-calc/internal/io"
)

var monthsInOneYear *big.Rat

func init() {
	monthsInOneYear = new(big.Rat)
	if _, ok := monthsInOneYear.SetString("12"); !ok {
		panic("set big rat to 12")
	}
}

type Bank struct {
	transactions  []io.Transaction
	interestRates []io.AnnualInterestRate
}

func NewBank(transactions []io.Transaction, interestRates []io.AnnualInterestRate) Bank {
	sort.SliceStable(transactions, func(i int, j int) bool {
		return transactions[i].Date.Before(transactions[j].Date)
	})

	sort.SliceStable(interestRates, func(i int, j int) bool {
		return interestRates[i].Day.Before(interestRates[j].Day)
	})

	return Bank{
		transactions:  transactions,
		interestRates: interestRates,
	}
}

// Process takes as input the state of a loan at the beginning of the
// given day and returns the state of the loan at the end of the same
// day.
func (b *Bank) Process(day time.Time, in Loan) (out Loan) {
	out = CopyLoan(in)

	if _, _, d := day.Date(); d == 1 {
		out.balance.Add(out.balance, out.interest)
		out.interest.Set(new(big.Rat))
	}

	trans := b.transactionsAmount(day)
	out.balance.Sub(out.balance, trans)

	rate, ok := b.annualInterestRate(day)
	if !ok {
		panic("annual interest rate not found")
	}

	y, m, _ := day.Date()
	dayRate := annualToDaily(rate, daysInMonth(m, y))
	dayInterest := new(big.Rat).Mul(dayRate, out.balance)

	out.interest.Add(out.interest, dayInterest)

	return out
}

func (b *Bank) transactionsAmount(day time.Time) *big.Rat {
	y, m, d := day.Date()
	day = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	amount := new(big.Rat)

	for _, t := range b.transactions {
		if ty, tm, td := t.Date.Date(); ty == y && tm == m && td == d {
			amount.Add(amount, t.Amount)
		} else if t.Date.After(day) {
			break
		}
	}

	return amount
}

func (b *Bank) annualInterestRate(day time.Time) (rate *big.Rat, ok bool) {
	day = DateFromTime(day)
	rate = new(big.Rat)

	for _, r := range b.interestRates {
		if r.Day.After(day) {
			return rate, ok
		} else {
			rate.Set(r.DecimalRate)
			ok = true
		}
	}

	return rate, ok
}

func daysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// annualToDaily calculates the daily interest rate for a month with a
// given number of days in it that corresponds to the specified annual
// rate.
func annualToDaily(rate *big.Rat, daysInMonth int) *big.Rat {
	daysInM := new(big.Rat).SetInt64(int64(daysInMonth))
	r := new(big.Rat).Set(rate)
	r.Quo(r, monthsInOneYear)
	r.Quo(r, daysInM)
	return r
}
