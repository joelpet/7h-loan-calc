package calc

import (
	"fmt"
	"math/big"
	"time"

	"gitlab.joelpet.se/joelpet/7h-loan-calc/internal/io"
)

// Run runs the calculations given the principal (i.e. initial sum of
// money borrowed), the first day of the loan, a list of transactions
// made, and a list of interest rate changes (incl. one that covers
// the first day of the loan).
func Run(firstDay time.Time, principal *big.Rat, interestRates []io.AnnualInterestRate, transactions []io.Transaction) {
	bank := NewBank(transactions, interestRates)
	loan := NewLoan(principal)

	start := DateFromTime(firstDay)
	end := DateFromTime(time.Now())

	fmt.Println("Date ; Annual interest rate (%) ; Balance ; Accrued interest")

	for day := DateFromTime(start); day.Before(end); day = day.AddDate(0, 0, 1) {
		loan = bank.Process(day, loan)

		air, ok := bank.annualInterestRate(day)

		airText := "-"
		if ok {
			airText = air.Mul(air, big.NewRat(100, 1)).FloatString(2)
		}

		fmt.Printf("%s ; %s ; %s ; %s \n",
			day.Format("2006-01-02"),
			airText,
			loan.balance.FloatString(2),
			loan.interest.FloatString(2))
	}
}

func DateFromTime(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
