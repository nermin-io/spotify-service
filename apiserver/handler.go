package apiserver

import (
	"go.uber.org/zap"
	"net/http"
)

func NewHandler(logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	return mux
}
