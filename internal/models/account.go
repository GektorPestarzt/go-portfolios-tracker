package models

type Status int

const (
	Created Status = iota
	Process
	Success
	Error
)

type Account struct {
	ID         int
	Token      string
	Type       string
	Username   string
	Status     Status
	Portfolios []*Portfolio
}

type Portfolio struct {
	TotalAmountPortfolio  *MoneyValue
	TotalAmountShares     *MoneyValue
	TotalAmountBonds      *MoneyValue
	TotalAmountEtf        *MoneyValue
	TotalAmountCurrencies *MoneyValue
	// TotalAmountFutures    MoneyValue
	ExpectedYield *Quotation
	Positions     []*Position
}

type Position struct {
	Figi                 string
	InstrumentType       string
	Quantity             *Quotation
	AveragePositionPrice *MoneyValue
	ExpectedYield        *Quotation
}

type MoneyValue struct {
	Currency string
	Units    int64
	Nano     int32
}

type Quotation struct {
	Units int64
	Nano  int32
}
