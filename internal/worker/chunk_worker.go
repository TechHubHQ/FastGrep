package worker

import (
	"fastgrep/internal/search"
	"os"
	"sync"
)

type ChunkWorkerPool struct {
	ChunkChan chan *search.SearchJob
	ErrorChan chan error
	LogFile *os.File
	ResultChan chan *search.Match
	WG *sync.WaitGroup
	MaxWorkers int
}

func NewChunkWorkerPool(maxWorkers int, logFile *os.File, resChan chan *search.Match) *ChunkWorkerPool {
	return &ChunkWorkerPool{
		WG:         &sync.WaitGroup{},
		ChunkChan:  make(chan *search.SearchJob, maxWorkers*2), // Buffered to keep workers busy
		ErrorChan:  make(chan error),
		ResultChan: resChan,
		MaxWorkers: maxWorkers,
		LogFile:    logFile,
	}
}

func (wp *ChunkWorkerPool) Start() {
	for workerID := range wp.MaxWorkers {
		wp.WG.Add(1)
		go wp.ChunkWorker(workerID)
	}
}

func (wp *ChunkWorkerPool) Stop() {
	close(wp.ChunkChan)
	wp.WG.Wait()
}

func (wp *ChunkWorkerPool) ChunkWorker(workerID int) {
	defer wp.WG.Done()
	for searchJob := range wp.ChunkChan {
		// fmt.Fprintf(wp.LogFile, "[ChunkWorker %d] Searching chunk (%d bytes)\n", workerID, len(searchJob.Data))
		searchJob.FindMatches()

		// Return the buffer to the pool
		if searchJob.Data != nil {
			PutBuffer(searchJob.Data)
		}
	}
}
