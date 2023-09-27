package response

// Ошибка валидации поля структуры.
type FieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// Ошибки валидации структуры, переданной в запросе.
type ErrorsResponse struct {
	FieldErrors []FieldError `json:"errors"`
}
