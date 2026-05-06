package search

import (
	"strings"
	"testing"
)

func TestFindMatches(t *testing.T) {
	tests := []struct {
		name         string
		searchString string
		data         string
		ignoreCase   bool
		wantMatches  int
	}{
		{
			name:         "Simple match",
			searchString: "hello",
			data:         "hello world\ngoodbye world",
			ignoreCase:   false,
			wantMatches:  1,
		},
		{
			name:         "No match",
			searchString: "missing",
			data:         "hello world\ngoodbye world",
			ignoreCase:   false,
			wantMatches:  0,
		},
		{
			name:         "Case insensitive match",
			searchString: "HELLO",
			data:         "hello world",
			ignoreCase:   true,
			wantMatches:  1,
		},
		{
			name:         "Case sensitive mismatch",
			searchString: "HELLO",
			data:         "hello world",
			ignoreCase:   false,
			wantMatches:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resChan := make(chan *Match, 10)
			sj := NewSearchJob("test.txt", tt.searchString, []byte(tt.data), 1, tt.ignoreCase, resChan)

			go func() {
				sj.FindMatches()
				close(resChan)
			}()

			count := 0
			for range resChan {
				count++
			}

			if count != tt.wantMatches {
				t.Errorf("FindMatches() got %d matches, want %d", count, tt.wantMatches)
			}
		})
	}
}

func BenchmarkFindMatches(b *testing.B) {
	data := []byte(strings.Repeat("This is a test line with some content\n", 1000)) // 1000 lines
	searchString := "content"
	resChan := make(chan *Match, 1000)

	for b.Loop() {
		sj := NewSearchJob("bench.txt", searchString, data, 1, false, resChan)
		go func() {
			sj.FindMatches()
		}()

		// Drain results
		for range 1000 {
			<-resChan
		}
	}
}

func BenchmarkFindMatchesLargeChunk(b *testing.B) {
	// 1MB of data
	line := "This is a reasonably long line of text that we will repeat many times to simulate a large chunk.\n"
	data := []byte(strings.Repeat(line, 1024*1024/len(line)))
	searchString := "repeat"
	numLines := len(data) / len(line)
	resChan := make(chan *Match, numLines+1)

	for b.Loop() {
		sj := NewSearchJob("bench.txt", searchString, data, 1, false, resChan)
		go func() {
			sj.FindMatches()
		}()

		// Drain results
		for range numLines {
			<-resChan
		}
	}
}
