package service

import (
	"fmt"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
	"github.com/anton-uvarenko/backend_school/internal/pkg/currency"
)

type CurrencyService struct {
	converter *currency.CurrencyConverter
}

func NewCurrencySevice(converter *currency.CurrencyConverter) *CurrencyService {
	return &CurrencyService{
		converter: converter,
	}
}

func (s CurrencyService) Rate() (float32, error) {
	rate, err := s.converter.GetUAHToUSD()
	if err != nil {
		fmt.Printf("%v: [%v]", pkg.ErrRate, err)
		return 0, err
	}

	return rate, nil
}
