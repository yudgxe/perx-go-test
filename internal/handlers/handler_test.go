package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/yudgxe/perx-go-test/internal/model"
	"github.com/yudgxe/perx-go-test/internal/queue"
	"github.com/yudgxe/perx-go-test/internal/response"
	"github.com/yudgxe/perx-go-test/internal/store"
)

// RandomOrderDiff сравнивает два массива и возврашает строку с различиями
func RandomOrderDiff[T any](got, want []T) string {
	var got_extra, want_extra []T
	for _, g := range got {
		var found bool
		for _, w := range want {
			if cmp.Equal(g, w) {
				found = true
				break
			}
		}
		if !found {
			got_extra = append(got_extra, g)
		}
	}
	for _, w := range want {
		var found bool
		for _, g := range got {
			if cmp.Equal(g, w) {
				found = true
				break
			}
		}
		if !found {
			want_extra = append(want_extra, w)
		}
	}
	var msg string
	for _, w := range want_extra {
		msg = fmt.Sprintf("- %+v\n%s", w, msg)
	}
	for _, g := range got_extra {
		msg = fmt.Sprintf("+ %+v\n%s", g, msg)
	}
	return msg
}

// testRequest вспомогательная функция для тестирования.
func testRequest[T any](t *testing.T, handler http.Handler, r *http.Request, status int, errors []response.FieldError, check func(*testing.T, T)) {
	t.Helper()
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	body, err := io.ReadAll(w.Result().Body)
	if err != nil {
		panic(fmt.Sprintf("Ошибка чтения тела ответа: %v", err))
	}
	defer w.Result().Body.Close()

	if got, want := w.Result().StatusCode, status; got != want {
		t.Fatalf("StatusCode: %v, ожидалось %v. Тело ответа: %s", got, want, body)
	}
	type errorsResponse struct {
		FieldErrors []response.FieldError `json:"errors"`
	}
	if check != nil {
		var resp T
		err := json.Unmarshal(body, &resp)
		if err != nil {
			t.Fatalf("Ошибка декодирования JSON в %T: %v. Тело ответа: %s", resp, err, body)
		}
		check(t, resp)
	}
	if len(errors) > 0 {
		var resp errorsResponse
		err := json.Unmarshal(body, &resp)
		if err != nil {
			t.Fatalf("Ошибка декодирования JSON в ErrorsResponse: %v. Тело ответа: %s", err, body)
		}

		if diff := RandomOrderDiff(resp.FieldErrors, errors); diff != "" {
			t.Fatalf("Отличаются ошибки в ответе (-ожидалось +получено):\n%s", diff)
		}
	}
}

func buildRequest(method string, url string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(fmt.Sprintf("Ошибка создания запроса: %v", err))
	}
	r.Header.Add("Content-type", "application/json")
	return r
}

func TestTaskCreate(t *testing.T) {
	s := store.NewTaskStore()
	q := queue.New(0)
	handler := http.StripPrefix("/api", NewHandler(q, s))
	body := `{"n" : 10, "d" : 10, "n1" : 10, "l" : 10, "ttl" : 10}`
	r := buildRequest(http.MethodPost, "/api/task/create", strings.NewReader(body))
	type queueError struct {
		Err string `json:"error"`
	}
	testRequest(t, handler, r, 400, nil, func(t *testing.T, r queueError) {
		if got, want := r.Err, queue.ErrQueueIsFull.Error(); got != want {
			t.Fatalf("Отличаются ошибки в ответе. Получена: %s, ожидалась: %s", got, want)
		}
	})

	q = queue.New(10)
	handler = http.StripPrefix("/api", NewHandler(q, s))

	for _, test := range []struct {
		name   string
		status int
		body   string
		errors []response.FieldError
	}{
		{
			name:   "valid",
			status: 200,
			body:   `{"n" : 10, "d" : 10, "n1" : 10, "l" : 10, "ttl" : 10}`,
			errors: nil,
		},
		{
			name:   "invalid",
			status: 422,
			body:   `{}`,
			errors: []response.FieldError{
				{Field: "n", Tag: "gt", Param: "0"},
				{Field: "ttl", Tag: "gt", Param: "0"},
				{Field: "l", Tag: "gt", Param: "0"},
				{Field: "d", Tag: "gt", Param: "0"},
				{Field: "n1", Tag: "required", Param: ""},
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := buildRequest(http.MethodPost, "/api/task/create", strings.NewReader(test.body))
			testRequest[any](t, handler, r, test.status, test.errors, nil)
		})
	}
}

func TestTaskList(t *testing.T) {
	q := queue.New(10)
	s := store.NewTaskStore()

	handler := http.StripPrefix("/api", NewHandler(q, s))

	t1 := s.Create(model.Task{})
	t2 := s.Create(model.Task{})

	type expectedErr struct {
		Err string `json:"error"`
	}

	for _, test := range []struct {
		name     string
		query    string
		status   int
		expected []model.Task
		errs     *expectedErr
	}{
		{
			name:     "no_query_valid",
			status:   200,
			query:    "",
			expected: []model.Task{t1, t2},
			errs:     nil,
		},
		{
			name:     "with_query_valid",
			status:   200,
			query:    "sorted=true",
			expected: []model.Task{t1, t2},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := buildRequest(http.MethodGet, "/api/task/list", nil)
			r.URL.RawQuery = test.query
			testRequest(t, handler, r, test.status, nil, func(t *testing.T, tasks []model.Task) {
				if diff := RandomOrderDiff(tasks, test.expected); diff != "" {
					t.Fatalf("Отличаются ошибки в ответе (-ожидалось +получено):\n%s", diff)
				}
				if test.query != "" {
					for i := 0; i < len(tasks); i++ {
						if got, want := tasks[i], test.expected[i]; !cmp.Equal(got, want) {
							t.Fatalf("Порядок объектов в ответет отличается. Ожидалось: %v, получено: %v, индекс %d", got, want, i)
						}
					}
				}
			})
		})
	}
}
