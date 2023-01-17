package calc

import (
	"bytes"
	"encoding/csv"
	"math/big"
	"path"
	"testing"
	"time"

	"gitlab.joelpet.se/joelpet/7h-loan-calc/internal/io"
)

func TestRun(t *testing.T) {
	transactions, err := io.ReadTransactions(
		path.Join("..", "testdata", "transactions.csv"),
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

	var outCSV bytes.Buffer
	csvReader := csv.NewReader(&outCSV)

	Run(&outCSV, firstDay, principal, interestRates, transactions, csvReader.Comma)

	records, err := csvReader.ReadAll()
	if err != nil {
		t.Errorf("reading CSV output: %s", err)
	}
	if want, got := 225, len(records); want != got {
		t.Errorf("wanted %d, but got %d", want, got)
	}

	wantDateBalances := map[string]string{
		"2022-08-09": "100206.04",
		"2022-08-10": "97202.14",
		"2022-08-11": "97202.14",
		"2022-10-26": "97496.24",
		"2022-10-27": "94396.24",
		"2022-10-28": "94396.24",
		"2022-11-02": "94609.64",
		"2022-11-03": "93309.64",
		"2022-11-04": "93309.64",
		"2022-11-21": "93309.64",
		"2022-11-22": "92173.64",
		"2022-11-23": "92173.64",
		"2022-12-31": "92378.36",
		"2023-01-01": "92581.59",
	}

	for _, record := range records {
		date := record[0]
		balance := record[2]

		if want, ok := wantDateBalances[date]; ok {
			t.Run(date, func(t *testing.T) {
				t.Logf("%+v", record)
				if want, got := want, balance; want != got {
					t.Errorf("wanted balance %s on date %s, but got %s", want, date, got)
				}
			})
		} else {
			t.Logf("%+v", record)
		}
	}
}
