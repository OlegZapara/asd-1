package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func printProgress(completed, total int) {
	if (completed*100)%total == 0 {
		percent := completed * 100 / total
		fmt.Printf("\033[1A\033[K")
		fmt.Printf("\rGenerating: %d/%d (%d.0%%)\n", completed, total, percent)
	}
}

func main() {
	sizeStr := os.Args[1]
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		log.Fatalf("Failed to convert %s to int: %v", sizeStr, err)
	}
	outfile := fmt.Sprintf("./files/in/%d.txt", size)
	file, err := os.Create(outfile)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 16*1024) // 16KB buffer
	defer writer.Flush()

	batchSize := 1000
	batch := make([]byte, 0, batchSize*12)

	fmt.Print("Generating: 0/100 (0.00%)\n")

	for i := 0; i < size; i++ {
		batch = append(batch, fmt.Sprintf("%d\n", rand.Int31())...)
		if len(batch) >= batchSize*12 {
			if _, err = writer.Write(batch); err != nil {
				log.Fatalf("Failed to write to file: %v", err)
			}
			batch = batch[:0]
		}
		printProgress(i+1, size)
	}

	if len(batch) > 0 {
		if _, err = writer.Write(batch); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
	}
}
