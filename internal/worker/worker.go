package worker

import (
	"log"
	"time"

	"github.com/yudgxe/perx-go-test/internal/logutil"
	"github.com/yudgxe/perx-go-test/internal/model"
	"github.com/yudgxe/perx-go-test/internal/queue"
	"github.com/yudgxe/perx-go-test/internal/store"
	"github.com/yudgxe/perx-go-test/internal/types"
)

type JobCallBack func(job interface{})

type Worker struct {
	channel *queue.Queue
	store   *store.TaskStore
}

func New(channel *queue.Queue, store *store.TaskStore) Worker {
	return Worker{
		channel: channel,
		store:   store,
	}
}

func (w *Worker) StartWorker() {
	go func() {
		log.Printf("Запуск воркеров")
		for {
			if work, ok := w.channel.Consume().(model.Task); ok {
				if logutil.V(2) {
					log.Printf("Начало работы над заданием - id: %v", work.ID)
				}

				w.store.Update(work.ID, store.UpdateTaskParams{
					Status:    model.GetTaskStatusPointer(model.TaskStatusInProcess),
					StartedAt: types.NewNullTime(time.Now()),
				})

				tiker := time.NewTicker(time.Duration(work.L * float32(time.Second)))
				defer tiker.Stop()

				result := work.N1
				for range tiker.C {
					work.Iteration++
					result += work.D
					w.store.Update(work.ID, store.UpdateTaskParams{
						Iteration: &work.Iteration,
						Result:    &result,
					})

					if work.Iteration == work.N {
						w.store.Update(work.ID, store.UpdateTaskParams{
							Status:  model.GetTaskStatusPointer(model.TaskStatusCompleted),
							EndedAt: types.NewNullTime(time.Now()),
							Result:  &result,
						})

						time.AfterFunc(time.Duration(work.TTL*float32(time.Second)), func() {
							w.store.Delete(work.ID)
							if logutil.V(2) {
								log.Printf("Удаление задания - id: %v", work.ID)
							}
						})

						break
					}
				}
			}
		}
	}()
}
