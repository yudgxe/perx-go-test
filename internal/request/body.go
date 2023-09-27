package request

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/yudgxe/perx-go-test/internal/logutil"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(fieldName)
}

// Body читает тело запроса в указанную структуру.
func Body(r *http.Request, obj interface{}) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return NewBadRequestErorr(http.StatusBadRequest, fmt.Errorf("Ошибка чтения тела запроса: %w", err))
	}
	if logutil.V(2) {
		log.Printf("Тело запроса: %s\n", b)
	}
	if err := json.Unmarshal(b, obj); err != nil {
		return NewBadRequestErorr(http.StatusBadRequest, fmt.Errorf("Ошибка декодирования JSON: %w", err))
	}
	return nil
}

// ValidateBody читает и валидирует тело запроса в указанную структуру
func ValidateBody(r *http.Request, obj interface{}) error {
	if err := Body(r, &obj); err != nil {
		return err
	}
	if err := validate.Struct(obj); err != nil {
		return err
	}
	return nil
}

// fieldName возвращает имя поля из тэга json вместо имени поля структуры.
func fieldName(f reflect.StructField) string {
	name := strings.SplitN(f.Tag.Get("json"), ",", 1)[0]
	if name == "-" {
		return ""
	}
	return name
}
