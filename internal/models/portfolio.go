package models

type Account struct {
	ID         int
	Token      string
	Type       string
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
	Positions     []*PortfolioPosition
}

type PortfolioPosition struct {
	Figi                 string
	InstrumentType       string
	Quantity             *Quotation
	AvaragePositionPrice *MoneyValue
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
