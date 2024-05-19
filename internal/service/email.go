package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
	"github.com/anton-uvarenko/backend_school/internal/pkg/currency"
	"github.com/anton-uvarenko/backend_school/internal/pkg/email"
	"github.com/jackc/pgx/v5/pgconn"
)

type EmailService struct {
	emailRepo         emailRepo
	emailSender       *email.EmailSender
	currencyConverter *currency.CurrencyConverter
}

func NewEmailService(emailRepo emailRepo, sender *email.EmailSender, converter *currency.CurrencyConverter) *EmailService {
	return &EmailService{
		emailSender:       sender,
		emailRepo:         emailRepo,
		currencyConverter: converter,
	}
}

type emailRepo interface {
	AddEmail(ctx context.Context, email string) error
	GetAll(ctx context.Context) ([]string, error)
}

func (s *EmailService) AddEmail(ctx context.Context, email string) error {
	err := s.emailRepo.AddEmail(ctx, email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			err := err.(*pgconn.PgError)
			if err.Code == "23505" {
				return pkg.ErrEmailConflict
			}
		}

		fmt.Printf("%v: [%v]\n", pkg.ErrDBInternal, err)
		return pkg.ErrDBInternal
	}

	return nil
}

func (s *EmailService) SendEmails(ctx context.Context) error {
	emails, err := s.emailRepo.GetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNoEmailsRegistered
		}

		fmt.Printf("%v: [%v]\n", pkg.ErrDBInternal, err)
		return pkg.ErrDBInternal
	}

	currentCurrency, err := s.currencyConverter.GetUAHToUSD()
	if err != nil {
		return err
	}
	for _, email := range emails {
		err = s.emailSender.SendEmail(email, fmt.Sprintf("current ratio uah to usd is %v", currentCurrency))
		if err != nil {
			fmt.Printf("%v: [%v]", pkg.ErrEmailSend, err)
		}
	}

	return nil
}
