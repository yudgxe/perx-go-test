package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/yudgxe/perx-go-test/internal/model"
	"github.com/yudgxe/perx-go-test/internal/queue"
	"github.com/yudgxe/perx-go-test/internal/request"
	"github.com/yudgxe/perx-go-test/internal/response"
	"github.com/yudgxe/perx-go-test/internal/store"
)

type RouterEnv struct {
	queue *queue.Queue     // Очередь
	store *store.TaskStore // Хранилише
}

func NewHandler(queue *queue.Queue, store *store.TaskStore) http.Handler {
	env := &RouterEnv{
		queue: queue,
		store: store,
	}

	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/task/create", env.taskCreate())
	router.HandlerFunc(http.MethodGet, "/task/list", env.taskList())

	return router
}

func (e *RouterEnv) taskCreate() http.HandlerFunc {
	type taskCreateBody struct {
		N   int     `json:"n"   validate:"gt=0"`
		D   float32 `json:"d"   validate:"gt=0"`
		N1  float32 `json:"n1"  validate:"required"`
		L   float32 `json:"l"   validate:"gt=0"`
		TTL float32 `json:"ttl" validate:"gt=0"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body taskCreateBody
		if err := request.ValidateBody(r, &body); err != nil {
			handleError(w, r, err)
			return
		}
		task := e.store.Create(model.Task{
			N:      body.N,
			D:      body.D,
			N1:     body.N1,
			L:      body.L,
			TTL:    body.TTL,
			Status: model.TaskStatusInQueue,
		})
		if err := e.queue.Publish(task); err != nil {
			handleError(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (e *RouterEnv) taskList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := e.store.ReadMap()
		for i, msg := range e.queue.Export() {
			if task, ok := msg.(model.Task); ok {
				task.QueueNumber = i + 1
				m[task.ID] = task
			}
		}
		tasks := make([]model.Task, 0, len(m))
		for _, task := range m {
			tasks = append(tasks, task)
		}
		query := r.URL.Query().Get("sorted")
		if query != "" {
			sorted, err := strconv.ParseBool(query)
			if err != nil {
				err = request.ErrInvalidBool
				handleError(w, r, err)
				return
			}
			if sorted {
				sort.SliceStable(tasks, func(i, j int) bool {
					return tasks[i].ID < tasks[j].ID
				})
			}
		}
		response.WriteJSON(w, http.StatusOK, tasks)
	}
}
