package main

import (
	"os"
	"fmt"
	"log"
	"sync"
	"time"
	"bufio"

	"github.com/VaticanEmbassy/golastipass/cfg"
	"github.com/VaticanEmbassy/golastipass/doc"
	"github.com/VaticanEmbassy/golastipass/util"
	"github.com/VaticanEmbassy/golastipass/elastic"
)

var wg sync.WaitGroup


func ProcessLine(writer *elastic.Writer, type_ string, source int,
		partitionLevels int, lines chan string) {
	for {
		line, ok := <- lines
		if ok {
			d, err := doc.ParseLine(line, type_, source, partitionLevels)
			if err == nil {
				writer.Add(&d)
			}
		} else {
			wg.Done()
			return
		}
	}
}


func ProcessFile(fname string, ch chan string) int {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	start := time.Now()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count += 1
		line := scanner.Text()
		ch <- line
		if (count % 100000 == 0) {
			fmt.Println(fmt.Sprintf("processed %d lines of file %s",
						count, fname))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	fmt.Println(fmt.Sprintf("finished processing %d lines of file %s in %s",
				count, fname, util.ElapsedTime(start, time.Now())))
	return count
}


func main() {
	config := cfg.ReadArgs()
	writer := elastic.NewWriter(config)
	start := time.Now()
	fmt.Print("creating indices, please wait... ")
	indicesNr := writer.CreateIndices()
	fmt.Printf("created %d indices in %s.\n", indicesNr,
		util.ElapsedTime(start, time.Now()))
	start = time.Now()

	ch := make(chan string, config.BufferedLines)
	for w := 1; w <= config.LineWorkers; w++ {
		wg.Add(1)
		go ProcessLine(writer, config.Type, config.Source,
				config.PartitionLevels, ch)
	}

	count := 0
	for _, fname := range config.Files {
		count += 1
		ProcessFile(fname, ch)
	}

	close(ch)
	wg.Wait()
	writer.Close()
	fmt.Printf("processed %d files in %v\n",
		count, util.ElapsedTime(start, time.Now()))
}
