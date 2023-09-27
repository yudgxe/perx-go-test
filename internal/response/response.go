package response

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/yudgxe/perx-go-test/internal/logutil"
)

func WriteValidationErrors(w http.ResponseWriter, validationErrors validator.ValidationErrors) {
	var response ErrorsResponse
	for _, err := range validationErrors {
		var fieldError FieldError
		fieldError.Field = err.Field()
		if ns := err.Namespace(); ns != "" {
			if idx := strings.Index(ns, "."); idx > 0 {
				ns = ns[idx+1:]
			}
			fieldError.Field = ns
		}
		fieldError.Tag = err.Tag()
		fieldError.Param = err.Param()
		response.FieldErrors = append(response.FieldErrors, fieldError)
	}
	WriteJSON(w, http.StatusUnprocessableEntity, response)
}

func WriteJSON(w http.ResponseWriter, status int, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if logutil.V(2) {
		log.Printf("Возвращаемое тело ответа: %s\n", b)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}
