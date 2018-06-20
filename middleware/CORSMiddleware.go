package middleware

import (
	"github.com/gorilla/handlers"
	"net/http"
)

/*
	CorsMiddleware is used for Allowing Cross-origin requests.
*/
func CorsMiddleware(h http.Handler, origins []string, methods []string, headers []string) http.Handler {
	originsOk := handlers.AllowedOrigins(origins)
	methodsOk := handlers.AllowedMethods(methods)
	headersOk := handlers.AllowedHeaders(headers)

	return handlers.CORS(originsOk, methodsOk, headersOk)(h)
}
