package calc

import (
	"math/big"
	"path"
	"testing"
	"time"

	"gitlab.joelpet.se/joelpet/7h-loan-calc/internal/io"
)

func TestRun(t *testing.T) {
	transactions, err := io.ReadTransactions(
		path.Join("..", "testdata", "transaktioner_20220810_20220810.csv"),
		';',
	)
	if err != nil {
		t.Errorf("reading transactions: %s", err)
	}
	t.Logf("Transactions:\n%+v", transactions)

	interestRates, err := io.ReadInterestRates(
		path.Join("..", "testdata", "annual_interest_rates.csv"),
		';',
	)
	if err != nil {
		t.Errorf("reading interest rates: %s", err)
	}
	t.Logf("Interest rates:\n%+v", interestRates)

	principal, ok := new(big.Rat).SetString("100_000")
	if !ok {
		t.Error("set big rat string")
	}

	firstDay := time.Date(2022, time.June, 7, 0, 0, 0, 0, time.UTC)

	Run(firstDay, principal, interestRates, transactions, ';')
}
