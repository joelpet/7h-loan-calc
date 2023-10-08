package calc

import (
	"encoding/csv"
	"io"
	"math/big"
	"time"

	intio "gitlab.joelpet.se/joelpet/7h-loan-calc/internal/io"
)

// Run runs the calculations given the principal (i.e. initial sum of
// money borrowed), the first day of the loan, a list of transactions
// made, and a list of interest rate changes (incl. one that covers
// the first day of the loan). Results are written to w as CSV
// records—one record per day—indicating the state of the loan on each
// day.
func Run(w io.Writer, firstDay time.Time, principal *big.Rat, interestRates []intio.AnnualInterestRate, transactions []intio.Transaction, outComma rune) {
	bank := NewBank(transactions, interestRates)
	loan := NewLoan(principal)

	start := DateFromTime(firstDay)
	end := DateFromTime(time.Now()).AddDate(0, 0, 1)

	writer := csv.NewWriter(w)
	writer.Comma = outComma

	defer writer.Flush()

	writer.Write([]string{
		"Date",
		"Annual interest rate (%)",
		"Balance",
		"Accrued interest",
	})

	for day := DateFromTime(start); day.Before(end); day = day.AddDate(0, 0, 1) {
		loan = bank.Process(day, loan)

		air, ok := bank.annualInterestRate(day)

		airText := "-"
		if ok {
			airText = air.Mul(air, big.NewRat(100, 1)).FloatString(2)
		}

		writer.Write([]string{
			day.Format("2006-01-02"),
			airText,
			loan.balance.FloatString(2),
			loan.interest.FloatString(2),
		})
	}
}

func DateFromTime(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
