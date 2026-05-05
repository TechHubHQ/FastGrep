package search

import (
	"bytes"
	"strings"
)

type Match struct {
	FileName string
	Content  []byte
	LineNum  int
}

func NewMatch(fileName string, lineNum int, content []byte) *Match {
	return &Match{
		FileName: fileName,
		LineNum:  lineNum,
		Content:  content,
	}
}

type SearchJob struct {
	ResultChan   chan *Match
	FileName     string
	SearchString string
	Data         []byte
	StartLine    int
	IgnoreCase   bool
}

func NewSearchJob(fileName, searchString string, chunk []byte, startLine int, ignoreCase bool, resChan chan *Match) *SearchJob {
	return &SearchJob{
		FileName:     fileName,
		SearchString: searchString,
		Data:         chunk,
		ResultChan:   resChan,
		StartLine:    startLine,
		IgnoreCase:   ignoreCase,
	}
}

func (sj *SearchJob) FindMatches() {
	lines := bytes.SplitSeq(sj.Data, []byte{'\n'})

	// Prepare patterns
	patternBytes := []byte(sj.SearchString)
	searchPattern := patternBytes

	// \033[33m = Yellow, \033[1m = Bold
	highlightedPattern := []byte("\033[33;1m" + sj.SearchString + "\033[0m")

	if sj.IgnoreCase {
		searchPattern = bytes.ToLower(patternBytes)
		highlightedPattern = []byte("\033[33;1m" + strings.ToLower(sj.SearchString) + "\033[0m")
	}

	lineOffset := 0
	for line := range lines {
		absoluteLineNum := sj.StartLine + lineOffset

		compareLine := line
		if sj.IgnoreCase {
			compareLine = bytes.ToLower(line)
		}

		if bytes.Contains(compareLine, searchPattern) {
			highlightedLine := bytes.ReplaceAll(compareLine, searchPattern, highlightedPattern)

			match := NewMatch(sj.FileName, absoluteLineNum, highlightedLine)
			sj.ResultChan <- match
		}
		lineOffset++
	}
}
