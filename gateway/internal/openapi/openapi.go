package openapi

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func RegisteraOpenAPIEndpoint(mux *runtime.ServeMux, dir string) error {
	openAPIdata, err := loadOpenAPIs(dir)
	if err != nil {
		return fmt.Errorf("app.openapi: %w", err)
	}

	mux.HandlePath(http.MethodGet, "/api.swagger.json", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(openAPIdata)
		w.WriteHeader(http.StatusOK)
	})

	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/api.swagger.json"),
	)
	err = mux.HandlePath(http.MethodGet, "/swagger/*", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		swaggerHandler(w, r)
	})
	if err != nil {
		return fmt.Errorf("app.openapi: %w", err)
	}
	return nil
}

// TODO: use single file
// loadOpenAPIs load open api files from direcotory
// and merge them into single open api file
func loadOpenAPIs(dir string) ([]byte, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}
