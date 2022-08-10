package main

import (
	"flag"
	"log"
	"math/big"
	"time"

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
)

func main() {
	log.SetFlags(0)

	flag.BoolVar(&version, "v", false, "print the version")
	flag.StringVar(&transactions, "t", "transactions.csv", "transactions CSV `file`")
	flag.StringVar(&interestRates, "r", "interest_rates.csv", "interest rates CSV `file`")
	flag.StringVar(&firstDay, "d", "2022-06-27", "`date` of first day of loan")
	flag.StringVar(&principal, "p", "200000", "principal `balance` on first day")

	flag.Parse()

	if version {
		log.Print(buildinfo.Version())
		return
	}

	firstDayT, err := time.Parse("2006-01-02", firstDay)
	if err != nil {
		log.Fatalf("failed to read first day argument: %s", err)
	}

	principalBalance, ok := new(big.Rat).SetString(principal)
	if !ok {
		log.Fatalf("failed to parse principal balance %q", principal)
	}

	transactionsL, err := intio.ReadTransactions(transactions)
	if err != nil {
		log.Fatalf("failed to read transactions: %s", err)
	}

	interestRatesL, err := intio.ReadInterestRates(interestRates)
	if err != nil {
		log.Fatalf("failed to read interest rates: %s", err)
	}

	log.Printf("Calculating loan based on %d transaction(s) and %d interest rate entries.",
		len(transactionsL), len(interestRatesL))

	calc.Run(firstDayT, principalBalance, interestRatesL, transactionsL)
}
