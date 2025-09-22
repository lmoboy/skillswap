package config

import (
	"net/http"

	"github.com/rs/cors"
)

func CORS() *cors.Cors {
	return cors.New(cors.Options{
		// Saraksts ar atļautajiem domēniem, no kuriem var veikt pieprasījumus.
		AllowedOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		// Atļauj sūtīt akreditācijas datus, piemēram, sīkdatnes (cookies).
		AllowCredentials: true,
		// Saraksts ar atļautajām HTTP metodēm (GET, POST, PUT, utt.).
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		// Saraksts ar atļautajām galvenēm (headers), ko klients var nosūtīt.
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
		},
		// Saraksts ar galvenēm, kuras serveris atļauj eksponēt klientam.
		ExposedHeaders: []string{
			"Content-Length",
			"Set-Cookie",
		},
		// Atspējo CORS atkļūdošanas (debug) izvadi konsolē.
		Debug: false,
	})
}
