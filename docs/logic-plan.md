# FastGrep Logic

collect the arguments from CLI and initialize a concurrent search function we assign single worker per file if multiple files

- consider worker-pool pattern
- worker pool pattern can init multiple workers on same task and assign each area to each worker like a single file to each

- again in the worker pool pattern we initialize set of goroutines to search the pattern across chunks of texts in a single file to allow faster search

- lets say we read 5 chunks with 5 goroutines once the 1st goroutine is completed it will go to 6th chunk to continue

- in this way we will create number of goroutines as per the files given and in that fixed number of chunk reading goroutines to process search
