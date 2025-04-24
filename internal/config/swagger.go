package config

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/fx"
)

// SwaggerModule provides the Swagger configuration as an fx module
var SwaggerModule = fx.Options(
	fx.Invoke(ConfigureSwagger),
)

// ConfigureSwagger sets up Swagger routes on an existing router
func ConfigureSwagger(router *chi.Mux, cfg *Config) {
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(cfg.ServerURL+"/swagger/doc.json"),
	))
}
