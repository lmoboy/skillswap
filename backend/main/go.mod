module GoBackend

go 1.25.0

replace service.com/services => ../services

require (
	github.com/gorilla/mux v1.8.1
	github.com/rs/cors v1.11.1
)
