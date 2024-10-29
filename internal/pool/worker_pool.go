package pool

import (
	"sync"
)

// WorkerPool представляет пул воркеров для выполнения задач.
type WorkerPool struct {
	taskQueue  chan func() // Очередь задач
	workerPool int         // Текущее количество воркеров
	wg         sync.WaitGroup
	mu         sync.Mutex    // Для защиты данных
	done       chan struct{} // Канал для завершения воркеров
}

// NewWorkerPool создает новый пул воркеров с заданным количеством воркеров.
func NewWorkerPool(workerCount int) *WorkerPool {
	wp := &WorkerPool{
		taskQueue:  make(chan func(), 100),
		workerPool: workerCount,
		done:       make(chan struct{}),
	}

	wp.StartWorkers(workerCount)
	return wp
}

// StartWorkers запускает указанное количество воркеров.
func (wp *WorkerPool) StartWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go wp.worker()
	}
}

// SetWorkerCount обновляет количество воркеров.
func (wp *WorkerPool) SetWorkerCount(newWorkerCount int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	// Останавливаем существующих воркеров.
	close(wp.done)

	// Создаем новый канал для завершения воркеров.
	wp.done = make(chan struct{})

	// Запускаем новых воркеров.
	wp.StartWorkers(newWorkerCount)
	wp.workerPool = newWorkerCount

}

// worker запускает задачи из очереди.
func (wp *WorkerPool) worker() {
	for {
		select {
		case task := <-wp.taskQueue:
			// Выполняем задачу.
			task()
			wp.wg.Done() // Уменьшаем счетчик после выполнения задачи
		case <-wp.done:
			// Завершение работы воркера.
			return
		}
	}
}

// SubmitTask добавляет задачу в очередь
func (wp *WorkerPool) SubmitTask(task func()) {
	wp.wg.Add(1)         // Увеличиваем счетчик перед добавлением задачи
	wp.taskQueue <- task // Добавляем задачу в очередь
}

// Wait завершает выполнение всех задач.
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Close завершает работу пула воркеров.
func (wp *WorkerPool) Close() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	close(wp.done)
	wp.Wait()
}
