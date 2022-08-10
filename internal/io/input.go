package io

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Transaction struct {
	// Date marks the year, month, and day when this transaction
	// happened in UTC. All other fields should be zero.
	Date        time.Time `csv:"Datum"`
	Type        string    `csv:"Typ av transaktion"`
	Description string    `csv:"Värdepapper/beskrivning"`
	Amount      *big.Rat  `csv:"Belopp"`
	Currency    string    `csv:"Valuta"`
}

func ReadTransactions(csvFilename string) ([]Transaction, error) {
	file, err := os.Open(csvFilename)
	if err != nil {
		return nil, fmt.Errorf("opening CSV file: %w", err)
	}
	defer file.Close()

	transactions := []Transaction{}

	r := csv.NewReader(file)
	r.Comma = ';'

	_, err = r.Read() // skip first line

	for {
		r, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading transaction CSV record: %w", err)
		}

		rDate, rType, rDesc, rAmount, rCurrency :=
			r[0], r[2], r[3], r[6], r[8]

		date, err := time.Parse("2006-01-02", rDate)
		if err != nil {
			return nil, fmt.Errorf("parsing date: %w", err)
		}

		amount, err := ParseAmount(rAmount)
		if err != nil {
			return nil, fmt.Errorf("parsing amount %q: %w", rAmount, err)
		}

		transactions = append(transactions, Transaction{
			Date:        date,
			Type:        rType,
			Description: rDesc,
			Amount:      amount,
			Currency:    rCurrency,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("reading all transactions: %w", err)
	}

	return transactions, nil
}

func ParseAmount(amount string) (*big.Rat, error) {
	amount = strings.Replace(amount, ",", ".", 1)
	amount = strings.ReplaceAll(amount, " ", "")

	f, ok := new(big.Rat).SetString(amount)
	if !ok {
		return nil, fmt.Errorf("setting big rat to %q", amount)
	}

	return f, nil
}

func newTransaction(year int, month time.Month, day int, amount string) (Transaction, error) {
	rnd := make([]byte, 16)
	if _, err := rand.Read(rnd); err != nil {
		return Transaction{}, fmt.Errorf("reading random data: %w", err)
	}

	amnt, ok := new(big.Rat).SetString(amount)
	if !ok {
		return Transaction{}, fmt.Errorf("setting amount to string %q", amount)
	}

	return Transaction{
		Date:        time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		Type:        "Insättning",
		Description: base64.StdEncoding.EncodeToString(rnd),
		Amount:      amnt,
		Currency:    "SEK",
	}, nil
}

func MustNewTransaction(year int, month time.Month, day int, amount string) Transaction {
	if t, err := newTransaction(year, month, day, amount); err != nil {
		panic(err)
	} else {
		return t
	}
}

type AnnualInterestRate struct {
	Day         time.Time
	DecimalRate *big.Rat
}

func (r AnnualInterestRate) Equal(s AnnualInterestRate) bool {
	return r.Day.Equal(s.Day) && r.DecimalRate.Cmp(s.DecimalRate) == 0
}

func (r AnnualInterestRate) String() string {
	return fmt.Sprintf("%v %s", r.Day.Format("2006-01-02"), r.DecimalRate)
}

func NewAnnualInterestRate(year int, month time.Month, day int, decimalRate string) (AnnualInterestRate, error) {
	rate, ok := new(big.Rat).SetString(decimalRate)
	if !ok {
		return AnnualInterestRate{}, fmt.Errorf("setting big rat from string %q", decimalRate)
	}
	return AnnualInterestRate{
		Day:         time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		DecimalRate: rate,
	}, nil
}

func MustNewAnnualInterestRate(year int, month time.Month, day int, decimalRate string) AnnualInterestRate {
	if rate, err := NewAnnualInterestRate(year, month, day, decimalRate); err != nil {
		panic(err)
	} else {
		return rate
	}
}

func ReadInterestRates(csvFilename string) ([]AnnualInterestRate, error) {
	file, err := os.Open(csvFilename)
	if err != nil {
		return nil, fmt.Errorf("opening CSV file: %w", err)
	}
	defer file.Close()

	rates := []AnnualInterestRate{}
	bigRat100 := big.NewRat(100, 1)

	r := csv.NewReader(file)
	r.Comma = ';'

	_, err = r.Read() // skip first line

	for {
		r, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading interest rate CSV record: %w", err)
		}

		rDate, rPercentage := r[0], r[1]

		date, err := time.Parse("2006-01-02", rDate)
		if err != nil {
			return nil, fmt.Errorf("parsing date: %w", err)
		}

		percentage, ok := new(big.Rat).SetString(rPercentage)
		if !ok {
			return nil, fmt.Errorf("setting big rat to %q", percentage)
		}

		rates = append(rates, AnnualInterestRate{
			Day:         date,
			DecimalRate: percentage.Quo(percentage, bigRat100),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("reading all interest rates: %w", err)
	}

	return rates, nil

}
