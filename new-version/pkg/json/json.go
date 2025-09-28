package json

import (
	"encoding/json"
	"log/slog"
	"net/http"
	help "new-version/pkg/http-helpers"
)

func ReadRequestBody(r *http.Request, result any) error {
	slog.Debug("request body: ", result)

	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return err
	}

	return nil
}

func WriteResponseBody(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func WriteSuccess(w http.ResponseWriter, msg string, data any, code int) {
	resp := help.NewResponse(msg, data, code)
	WriteResponseBody(w, resp, code)
}

func WriteError(w http.ResponseWriter, msg string, code int) {
	resp := help.NewErrResponse(msg, code)
	WriteResponseBody(w, resp, code)
}
