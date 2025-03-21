package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorResponse represents the standard matrix API error response
type ErrorResponse struct {
	ErrCode string `json:"errcode"`
	ErrMsg  string `json:"error"`
}

// Implement type Error interface
func (e ErrorResponse) Error() string {
	return e.ErrCode + ": " + e.ErrMsg
}

// WriteErrorResponse sends Matrix API error response and logs the error.
func WriteErrorResponse(w http.ResponseWriter, statusCode int, errCode, errMsg string) {
	err := ErrorResponse{
		ErrCode: errCode,
		ErrMsg:  errMsg,
	}

	slog.Error(
		"API error",
		slog.String("errcode", err.ErrCode),
		slog.String("error", err.ErrMsg),
		slog.Int("status_code", statusCode),
	)

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}
