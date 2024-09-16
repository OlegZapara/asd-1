package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

type FileChunkSorter struct {
	inputFile   string
	outputFile  string
	chunkSize   int
	tempFiles   []string
	totalLines  int
	linesSorted int
}

func NewFileChunkSorter(inputFile, outputFile string, chunkSize int) *FileChunkSorter {
	return &FileChunkSorter{
		inputFile:  inputFile,
		outputFile: outputFile,
		chunkSize:  chunkSize,
		tempFiles:  []string{},
	}
}

func (f *FileChunkSorter) countTotalLines() (int, error) {
	file, err := os.Open(f.inputFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalLines := 0
	for scanner.Scan() {
		totalLines++
	}
	return totalLines, scanner.Err()
}

func (f *FileChunkSorter) SortFile() error {
	totalLines, err := f.countTotalLines()
	if err != nil {
		return err
	}
	f.totalLines = totalLines

	err = f.splitAndSortChunks()
	if err != nil {
		return err
	}
	return f.mergeChunks()
}

func (f *FileChunkSorter) splitAndSortChunks() error {
	file, err := os.Open(f.inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	chunk := []int{}

	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return err
		}
		chunk = append(chunk, num)

		f.linesSorted++
		f.printProgress("Splitting & Sorting", f.linesSorted, f.totalLines)

		if len(chunk) >= f.chunkSize {
			err = f.sortAndWriteChunk(chunk)
			if err != nil {
				return err
			}
			chunk = []int{}
		}
	}

	if len(chunk) > 0 {
		err = f.sortAndWriteChunk(chunk)
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (f *FileChunkSorter) sortAndWriteChunk(chunk []int) error {
	sort.Ints(chunk)

	tempFile, err := os.CreateTemp("", "chunk_")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	for _, num := range chunk {
		_, err := fmt.Fprintln(tempFile, num)
		if err != nil {
			return err
		}
	}

	f.tempFiles = append(f.tempFiles, tempFile.Name())
	return nil
}

func (f *FileChunkSorter) mergeChunks() error {
	var files []*os.File
	for _, tempFileName := range f.tempFiles {
		file, err := os.Open(tempFileName)
		if err != nil {
			return err
		}
		files = append(files, file)
		defer file.Close()
		defer os.Remove(tempFileName)
	}

	output, err := os.Create(f.outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	readers := make([]*bufio.Scanner, len(files))
	for i, file := range files {
		readers[i] = bufio.NewScanner(file)
		readers[i].Scan()
	}

	linesWritten := 0

	for {
		minIndex := -1
		for i, reader := range readers {
			if reader == nil {
				continue
			}
			if minIndex == -1 || compareScanners(readers[minIndex], reader) > 0 {
				minIndex = i
			}
		}

		if minIndex == -1 {
			break
		}

		_, err := fmt.Fprintln(output, readers[minIndex].Text())
		if err != nil {
			return err
		}

		linesWritten++
		f.printProgress("Merging", linesWritten, f.totalLines)

		if !readers[minIndex].Scan() {
			readers[minIndex] = nil
		}
	}

	return nil
}

func compareScanners(s1, s2 *bufio.Scanner) int {
	num1, _ := strconv.Atoi(s1.Text())
	num2, _ := strconv.Atoi(s2.Text())
	return num1 - num2
}

func (f *FileChunkSorter) printProgress(stage string, completed, total int) {
	if (completed*100)%total == 0 {
		percent := (float64(completed) / float64(total)) * 100
		fmt.Printf("\033[1A\033[K")
		fmt.Printf("\r%s: %d/%d (%.2f%%)\n", stage, completed, total, percent)
	}
}

func main() {
	fmt.Println("Started sorting...")
	start := time.Now()
	sizeStr := os.Args[1]
	_, err := strconv.Atoi(sizeStr)
	if err != nil {
		fmt.Println("Failed to convert", sizeStr, "to int:", err)
		return
	}
	inputFile := "./files/in/" + sizeStr + ".txt"
	outputFile := "./files/out/" + sizeStr + ".txt"
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}

	fileStat, _ := file.Stat()
	chunkSize := int(fileStat.Size() / 100)
	fmt.Println("Chunk size:", chunkSize)

	sorter := NewFileChunkSorter(inputFile, outputFile, chunkSize)
	err = sorter.SortFile()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Sorting completed. Output written to", outputFile)
	}

	elapsed := time.Since(start)
	fmt.Println("Time taken:", elapsed)
}
