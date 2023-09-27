package model

import (
	"time"

	"github.com/yudgxe/perx-go-test/internal/types"
)

type TaskStatus string

const (
	TaskStatusInProcess = "processing"
	TaskStatusInQueue   = "wait"
	TaskStatusCompleted = "completed"
)

func GetTaskStatusPointer(status TaskStatus) *TaskStatus {
	return &status
}

type Task struct {
	ID          int            `json:"id"`           // Идентификатор задачи
	N           int            `json:"n"`            // Количество элементов
	D           float32        `json:"d"`            // Дельта между элементами последовательности
	N1          float32        `json:"n1"`           // Стартовое значение
	L           float32        `json:"l"`            // Интервал в секундах между итерациями
	TTL         float32        `json:"ttl"`          // Время хранения результата в секундах
	Iteration   int            `json:"iteration"`    // Текущая итерация
	CreatedAt   time.Time      `json:"created_at"`   // Время постановки задачи
	StartedAt   types.NullTime `json:"started_at"`   // Время старта задачи
	EndedAt     types.NullTime `json:"ended_at"`     // Время окончания задачи
	Status      TaskStatus     `json:"status"`       // Состояние задачи
	Result      float32        `json:"result"`       // Результат задачи
	QueueNumber int            `json:"queue_number"` // Номер в очереди
}
