package gravy

// Exchange is the interface to an exchange, e.g. a mock exchange or an alpaca client.
type Exchange interface {
	//CurrentPortfolio() trading.Portfolio
	//SubmitOrder(order trading.Order) trading.OrderOutcome
}
