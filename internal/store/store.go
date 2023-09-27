package store

import (
	"sync"
	"time"

	"github.com/yudgxe/perx-go-test/internal/model"
	"github.com/yudgxe/perx-go-test/internal/types"
)

type TaskStore struct {
	mx       sync.RWMutex
	database map[int]model.Task

	currentIndex int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		database: make(map[int]model.Task),
	}
}

func (s *TaskStore) Create(task model.Task) model.Task {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.currentIndex++
	task.ID = s.currentIndex
	task.CreatedAt = time.Now()
	s.database[task.ID] = task
	return task
}

func (s * TaskStore) ReadMap() map[int]model.Task {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.database
}

type UpdateTaskParams struct {
	Status    *model.TaskStatus
	Result    *float32
	StartedAt *types.NullTime
	EndedAt   *types.NullTime
	Iteration *int
}

// Update обновляет поля если они не null
func (s *TaskStore) Update(id int, params UpdateTaskParams) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if task, ok := s.database[id]; ok {
		if params.Status != nil {
			task.Status = *params.Status
		}
		if params.Result != nil {
			task.Result = *params.Result
		}
		if params.StartedAt != nil {
			task.StartedAt = *params.StartedAt
		}
		if params.EndedAt != nil {
			task.EndedAt = *params.EndedAt
		}
		if params.Iteration != nil {
			task.Iteration = *params.Iteration
		}
		s.database[id] = task
	}
}

func (s *TaskStore) UpdateStatus(id int, status model.TaskStatus) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if task, ok := s.database[id]; ok {
		task.Status = status
		s.database[id] = task
	}
}

func (s *TaskStore) Delete(id int) {
	s.mx.Lock()
	delete(s.database, id)
	s.mx.Unlock()
}
