package apiserver

import (
	"github.com/nermin-io/spotify-service/apiserver/middleware"
	"go.uber.org/zap"
	"net/http"
)

func NewHandler(logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	return middleware.Apply(mux, middleware.NewLoggingMiddleware(logger))
}
