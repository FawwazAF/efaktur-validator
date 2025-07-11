package main

import (
	efaktur_controller "github.com/efaktur-validator/internal/controller/efaktur"
	"github.com/efaktur-validator/internal/repository/api"
	"github.com/efaktur-validator/internal/server/http"
	efaktur_handler "github.com/efaktur-validator/internal/server/http/efaktur"
	"github.com/efaktur-validator/internal/server/http/index"
)

type Application struct {
	HTTPServers *http.Server
}

func main() {
	app := new(Application)
	// ===================================================== REPOSITORY =====================================================
	apiRepo := api.New()
	// ===================================================== CONTROLLER =====================================================
	efakturController := efaktur_controller.New(apiRepo)
	// ===================================================== HANDLER ========================================================
	efakturHandler := efaktur_handler.New(efakturController)

	handler := http.Handler{
		Index:   index.NewHandler(),
		Efaktur: efakturHandler,
	}

	// Add Router
	app.HTTPServers = http.NewServer(handler)
	app.HTTPServers.Start(":8080")
}
