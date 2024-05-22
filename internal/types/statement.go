package types

type AccountStatement struct {
	Title               string `json:"title"`
	GeneratedDate       string `json:"generated_date"`
	BankName            string `json:"bank_name"`
	CustomerName        string `json:"customer_name"`
	CustomerAddress     `json:"customer_address"`
	BalanceSummary      []BalanceSummary      `json:"balance_summary"`
	AccountTransactions []AccountTransactions `json:"account_transactions"`
	IBANDetails         []IBANDetails         `json:"iban_details"`
}

func (as AccountStatement) GetTotalOpeningBalanceSummary() float64 {
	var total float64 = 0
	for _, bs := range as.BalanceSummary {
		total += bs.OpeningBalance
	}

	return total
}

func (as AccountStatement) GetTotalMoneyOutBalanceSummary() float64 {
	var total float64 = 0
	for _, bs := range as.BalanceSummary {
		total += bs.MoneyOut
	}

	return total
}

func (as AccountStatement) GetTotalMoneyInBalanceSummary() float64 {
	var total float64 = 0
	for _, bs := range as.BalanceSummary {
		total += bs.MoneyIn
	}

	return total
}

func (as AccountStatement) GetTotalClosingBalanceSummary() float64 {
	var total float64 = 0
	for _, bs := range as.BalanceSummary {
		total += bs.ClosingBalance
	}

	return total
}

type CustomerAddress struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	County       string `json:"county"`
	Postcode     string `json:"post_code"`
}

type BalanceSummary struct {
	Product        string  `json:"product"`
	OpeningBalance float64 `json:"opening_balance"`
	MoneyIn        float64 `json:"money_in"`
	MoneyOut       float64 `json:"money_out"`
	ClosingBalance float64 `json:"closing_balance"`
}

type AccountTransactions struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	MoneyIn     float64 `json:"money_in"`
	MoneyOut    float64 `json:"money_out"`
	Balance     float64 `json:"balance"`
}

type IBANDetails struct {
	IBAN string `json:"iban"`
	BIC  string `json:"bic"`
	Note string `json:"note"`
}
