package calc

// TODO: Test amortization done on the 1st (or 3rd -- whenever the rollover happens)

// TODO: Test what happens if interest is paid a day or two too early
// This should be gracefully handled by the dumping the interest that is due onto the loan.
// That way, a premature interest payment would end up giving a little less daily interest for a
// couple of days, and then when the rollover happens, there is nothing still "due".
