package handlers

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/yudgxe/perx-go-test/internal/queue"
	"github.com/yudgxe/perx-go-test/internal/request"
	"github.com/yudgxe/perx-go-test/internal/response"
)

// handleError обрабатывает ошибки записывая ответ на запрос.
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	var badRequestError request.BadRequestError
	if errors.As(err, &badRequestError) {
		http.Error(w, badRequestError.Error(), badRequestError.Status)
		return
	}
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		response.WriteValidationErrors(w, validationErrors)
		return
	}
	if errors.Is(err, queue.ErrQueueIsFull)  {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if errors.Is(err, request.ErrInvalidBool) {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
}
