package main

import (
	"net/http"
	"os"

	"github.com/anton-uvarenko/backend_school/internal/core"
	"github.com/anton-uvarenko/backend_school/internal/db"
	"github.com/anton-uvarenko/backend_school/internal/pkg/currency"
	"github.com/anton-uvarenko/backend_school/internal/pkg/email"
	"github.com/anton-uvarenko/backend_school/internal/pkg/server"
	"github.com/anton-uvarenko/backend_school/internal/service"
	"github.com/anton-uvarenko/backend_school/internal/transport"
	"github.com/go-co-op/gocron/v2"
)

func main() {
	conn := db.Connect()

	queries := core.New(conn)

	emailSender := email.NewEmailSender(os.Getenv("FROM_EMAIL"), os.Getenv("FROM_EMAIL_PASSWORD"))
	converter := currency.NewCurrencyConverter(http.DefaultClient)

	service := service.NewService(queries, emailSender, converter)
	handler := transport.NewHandler(service)

	httpServer := server.NewServer(handler)

	scheduler, _ := gocron.NewScheduler()
	server.RegisterJobs(scheduler, service.EmailService)
	scheduler.Start()

	httpServer.ListenAndServe()
}
