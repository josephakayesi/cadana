package domain

type Rate struct {
	Model
	CurrencyPair string `json:"currency_pair" gorm:"not null"`
	Value        string `json:"value" gorm:"not null"`
}
