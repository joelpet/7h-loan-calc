package calc

import (
	"math/big"
)

type Loan struct {
	// balance is the unpaid principal balance (UPB) of the loan.
	balance *big.Rat
	// interest is the accrued interest over some period of time –
	// typically a calendar month.
	interest *big.Rat
}

func NewLoan(principalBalance *big.Rat) Loan {
	return Loan{
		balance:  new(big.Rat).Set(principalBalance),
		interest: new(big.Rat),
	}
}

func CopyLoan(loan Loan) Loan {
	cpy := NewLoan(loan.balance)
	cpy.interest.Set(loan.interest)
	return cpy
}
