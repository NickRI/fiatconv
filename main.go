package main

import (
	"context"
	"os"

	"github.com/NickRI/fiatconv/converter/interfaces/presenters"

	"github.com/NickRI/fiatconv/converter/domain/usecases"
	"github.com/NickRI/fiatconv/converter/interfaces/repositories"
	"github.com/NickRI/fiatconv/external"

	"github.com/NickRI/fiatconv/converter/interfaces/controllers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	remoteRepo := repositories.NewRestapiRepo(external.NewExchangeRatesClient())
	intrerator := usecases.NewBaseInteractor(remoteRepo)
	presenter := presenters.NewCli(os.Stdout)

	console := controllers.NewCli(intrerator, presenter)

	if err := console.Convert(ctx, os.Args[1:]); err != nil {
		panic(err)
	}

	cancel()
}
