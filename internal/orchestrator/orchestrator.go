package orchestrator

import (
	"fastgrep/internal/fsutil"
	"fastgrep/internal/search"
	"fastgrep/internal/worker"
	"fmt"
	"os"
	"runtime"
	"sync"
)

func Execute(searchString string, files []string, ignoreCase bool, recursive bool) error {
	logFile, err := os.OpenFile("fastgrep.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open log file: %v", err)
	}
	defer logFile.Close()

	matchedFiles, err := fsutil.GetGlobs(files)
	if err != nil {
		return err
	}

	numValidFiles, validFiles := fsutil.FilterFiles(matchedFiles, recursive)
	if len(validFiles) == 0 {
		return fmt.Errorf("no valid files found to search")
	}

	// Create a shared result channel for all workers
	resChan := make(chan *search.Match, 1000)
	var printerWG sync.WaitGroup

	// Start the printer goroutine
	printerWG.Go(func() {
		for match := range resChan {
			// \033[34m = Blue, \033[31m = Red, \033[0m = Reset
			fmt.Printf("\033[34m%s\033[0m:\033[31m%d\033[0m:%s\n", match.FileName, match.LineNum, string(match.Content))
		}
	})

	// Initialize the shared ChunkWorkerPool
	numCPU := runtime.NumCPU()
	cwPool := worker.NewChunkWorkerPool(numCPU, logFile, resChan)
	cwPool.Start()

	numWorkers := min(numValidFiles, numCPU)
	fwPool := worker.NewFileWorkerPool(numWorkers, searchString, logFile, ignoreCase, resChan, cwPool)

	fmt.Fprintf(logFile, "\n--- NEW SEARCH START: %s ---\n", searchString)
	fmt.Fprintf(logFile, "[Orchestrator] Starting %d workers...\n", numWorkers)
	fwPool.WorkerPoolExecutor()

	// Init Error listener -> for errors coming from fwPool
	go func() {
		for err := range fwPool.ErrorChan {
			fmt.Fprintf(logFile, "[ERROR] %s\n", err.Error())
		}
	}()

	for _, file := range validFiles {
		fwPool.FilesChan <- file
	}
	close(fwPool.FilesChan)

	// Wait for file workers to finish dispatching jobs
	fwPool.WG.Wait()
	close(fwPool.ErrorChan)

	// Stop the chunk workers
	cwPool.Stop()

	// Close result channel and wait for printer to finish
	close(resChan)
	printerWG.Wait()

	fmt.Fprintf(logFile, "[Orchestrator] All work complete!\n")

	return nil
}
