package search

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
)

//Resut - discribe the result of one search
type Result struct {
	// phrase were looking for
	Phrase string
	//whole line in which the entry was found (without \n or \r\n)
	Line string
	//line number(starting from 1) in which the entry was found
	LineNum int64
	//position(column) number (starting form 1) in which the entry was found
	ColNum int64
}

//All - searches all phrase entries in files (text files).
func All(ctx context.Context, phrase string, files []string) <-chan []Result {

	ch := make(chan []Result)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)

	for i := 0; i < len(files); i++ {
		wg.Add(1)

		go func(ctx context.Context, path string, ch chan<- []Result) {
			defer wg.Done()

			results := []Result{}
			data, err := os.ReadFile(path)
			if err != nil {
				log.Print("can not open the file:", err)
			}

			dataStr := string(data)
			splitData := strings.Split(dataStr, "\n")

			for index, line := range splitData {
				if strings.Contains(line, phrase) {
					result := Result{
						Phrase:  phrase,
						Line:    line,
						LineNum: int64(index + 1),
						ColNum:  int64(strings.Index(line, phrase) + 1),
					}
					results = append(results, result)
				}
			}

			if len(results) > 0 {
				ch <- results
			}
		}(ctx, files[i], ch)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	cancel()
	return ch
}

//Any - searches anyone phrase entries in files (text files).
func Any(ctx context.Context, phrase string, files []string) <-chan Result {

	ch := make(chan Result)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)

	for i := 0; i < len(files); i++ {
		wg.Add(1)

		go func(ctx context.Context, path string, ch chan<- Result) {
			defer wg.Done()
			
			data, err := os.ReadFile(path)
			if err != nil {
				log.Print("can not open the file:", err)
			}

			dataStr := string(data)
			splitData := strings.Split(dataStr, "\n")

			for index, line := range splitData {
				select {
				case <-ctx.Done():
					return
				default:
					if strings.Contains(line, phrase) {
						result := Result{
							Phrase:  phrase,
							Line:    line,
							LineNum: int64(index + 1),
							ColNum:  int64(strings.Index(line, phrase) + 1),
						}
						ch <- result
						cancel()
					}
				}
			}
		}(ctx, files[i], ch)
	}

	go func() {
		defer close(ch)
		wg.Wait()
		cancel()
	}()

	return ch
}
