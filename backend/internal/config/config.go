package config

import (
	"net/http"

	"github.com/rs/cors"
)

func CORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			if origin == "localhost:3000" {
				return true
			}
			return false
		},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*",
		},
		ExposedHeaders: []string{
			"*",
		},
		Debug: false,
	})
}
