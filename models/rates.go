package models

type Income struct {
	Symbol     string  `json:"symbol" gorm:"primaryKey"`
	Price_24h  float64 `json:"price_24h"`
	Volume_24h float64 `json:"volume_24h"`
	Last_trade float64 `json:"last_trade_price"`
}
