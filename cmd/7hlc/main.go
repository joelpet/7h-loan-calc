package main

import (
	"errors"
	"flag"
	"log"
	"math/big"
	"os"
	"time"

	"gitlab.joelpet.se/joelpet/7h-loan-calc/internal"
	"gitlab.joelpet.se/joelpet/7h-loan-calc/internal/buildinfo"
	"gitlab.joelpet.se/joelpet/7h-loan-calc/internal/calc"
	intio "gitlab.joelpet.se/joelpet/7h-loan-calc/internal/io"
)

var (
	version       bool   // -v flag
	transactions  string // -t flag
	interestRates string // -r flag
	firstDay      string // -d flag
	principal     string // -p flag
	csvInComma    string // -n flag
	csvOutComma   string // -u flag
)

func checkCSVComma(csvComma string) (rune, error) {
	comma := []rune(csvComma)
	if len := len(comma); len != 1 {
		return rune(0), errors.New("must be a single character")
	} else {
		return comma[0], nil
	}
}

func main() {
	log.SetFlags(0)

	flag.BoolVar(&version, "v", false, "print the version")
	flag.StringVar(&transactions, "t", "transactions.csv", "transactions CSV `file`")
	flag.StringVar(&interestRates, "r", "interest_rates.csv", "interest rates CSV `file`")
	flag.StringVar(&firstDay, "d", "2022-06-27", "`date` of first day of loan")
	flag.StringVar(&principal, "p", "200000", "principal `balance` on first day")
	flag.StringVar(&csvInComma, "n", ";", "input CSV file field delimiter `character` ")
	flag.StringVar(&csvOutComma, "u", ";", "output CSV file field delimiter `character` ")

	flag.Parse()

	if version {
		log.Print(buildinfo.Version())
		return
	}

	inComma, err := checkCSVComma(csvInComma)
	if err != nil {
		log.Fatalf("failed to get input CSV file field delimiter character: %s", err)
	}

	outComma, err := checkCSVComma(csvOutComma)
	if err != nil {
		log.Fatalf("failed to get output CSV file field delimiter character: %s", err)
	}

	firstDayT, err := time.Parse(internal.DateLayout, firstDay)
	if err != nil {
		log.Fatalf("failed to read first day argument: %s", err)
	}

	principalBalance, ok := new(big.Rat).SetString(principal)
	if !ok {
		log.Fatalf("failed to parse principal balance %q", principal)
	}

	transactionsL, err := intio.ReadTransactions(transactions, inComma)
	if err != nil {
		log.Fatalf("failed to read transactions: %s", err)
	}

	interestRatesL, err := intio.ReadInterestRates(interestRates, inComma)
	if err != nil {
		log.Fatalf("failed to read interest rates: %s", err)
	}

	log.Printf("Calculating loan based on %d transaction(s) and %d interest rate entries.",
		len(transactionsL), len(interestRatesL))

	calc.Run(os.Stdout, firstDayT, principalBalance, interestRatesL, transactionsL, outComma)
}
