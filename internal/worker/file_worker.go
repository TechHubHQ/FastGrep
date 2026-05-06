package worker

import (
	"bufio"
	"bytes"
	"errors"
	"fastgrep/internal/config"
	"fastgrep/internal/fsutil"
	"fastgrep/internal/search"
	"fmt"
	"io"
	"os"
	"sync"
)

type FileWorkerPool struct {
	ChunkWorkerPool *ChunkWorkerPool
	ErrorChan       chan error
	FilesChan       chan string
	LogFile         *os.File
	ResultChan      chan *search.Match
	WG              *sync.WaitGroup
	SearchString    string
	MaxWorkers      int
	IgnoreCase      bool
}

func NewFileWorkerPool(maxWorkers int, searchString string, logFile *os.File, ignoreCase bool, resChan chan *search.Match, cwPool *ChunkWorkerPool) *FileWorkerPool {
	return &FileWorkerPool{
		WG:              &sync.WaitGroup{},
		IgnoreCase:      ignoreCase,
		SearchString:    searchString,
		FilesChan:       make(chan string),
		ErrorChan:       make(chan error),
		MaxWorkers:      maxWorkers,
		LogFile:         logFile,
		ResultChan:      resChan,
		ChunkWorkerPool: cwPool,
	}
}

func (wp *FileWorkerPool) WorkerPoolExecutor() {
	for workerID := range wp.MaxWorkers {
		wp.WG.Add(1)
		go wp.FileWorker(workerID)
	}
}

func (wp *FileWorkerPool) FileWorker(workerID int) {
	defer wp.WG.Done()
	fmt.Fprintf(wp.LogFile, "[FileWorker %d] Started and waiting for files...\n", workerID)
	for file := range wp.FilesChan {
		fileSize, err := fsutil.GetFileSize(file)
		if err != nil {
			wp.ErrorChan <- err
			continue
		}

		numChunks := (fileSize + config.DefaultChunkSize - 1) / config.DefaultChunkSize
		fileObj, err := os.Open(file)
		if err != nil {
			wp.ErrorChan <- err
			continue
		}

		reader := bufio.NewReader(fileObj)
		var lineNum = 1
		for range numChunks {
			// Acquire buffer from pool
			chunk := GetBuffer()

			numReadBytes, err := io.ReadFull(reader, chunk)
			if err == io.EOF {
				PutBuffer(chunk)
				break
			}

			remainder, remErr := reader.ReadBytes('\n')
			if remErr != nil && remErr != io.EOF {
				wp.ErrorChan <- remErr
			}

			// Combine chunk and remainder into a fresh allocation to be safe.
			// This ensures the pooled 'chunk' can be returned immediately.
			fullLineChunk := make([]byte, numReadBytes+len(remainder))
			copy(fullLineChunk, chunk[:numReadBytes])
			copy(fullLineChunk[numReadBytes:], remainder)

			// Return the original 1MB buffer to the pool immediately
			PutBuffer(chunk)

			sJob := search.NewSearchJob(file, wp.SearchString, fullLineChunk, lineNum, wp.IgnoreCase, wp.ResultChan)
			wp.ChunkWorkerPool.ChunkChan <- sJob

			lineNum += bytes.Count(fullLineChunk, []byte{'\n'})
			if remErr == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
		}

		fileObj.Close()
		fmt.Fprintf(wp.LogFile, "[FileWorker %d] Finished file: %s\n", workerID, file)
	}
}
