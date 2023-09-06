package model

type Pendencia struct {
	Nome        string  `json:"nome"`
	Dia         string  `json:"dia"`
	Description string  `json:"description"`
	Valor       float64 `json:"valor"`
	Pago        bool    `json:"pago"`
}
