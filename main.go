package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"sync"
	"time"
	"bufio"
	"strconv"
	"io/ioutil"
	"path/filepath"

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
	var file *os.File
	var err error
	if fname == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(fname)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	start := time.Now()
	reader := bufio.NewReader(file)
	count := 0
	saveProgress := false
	linesToSkip := 0
	if fname != "-" {
		saveProgress = true
		linesToSkip = getProcessedLines(fname)
	}
	if linesToSkip != 0 {
		fmt.Printf("skipping %d lines that were already processed\n", linesToSkip)
		skipped := 0
		_, err = reader.ReadString('\n')
		for {
			skipped += 1
			if skipped >= linesToSkip {
				count = skipped
				break
			}
			if err == io.EOF {
				break
			}
			_, err = reader.ReadString('\n')
		}
	}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		count += 1
		ch <- line
		if (saveProgress && count % 10000 == 0) {
			storeProgress(fname, count)
		}
		if (count % 100000 == 0) {
			fmt.Println(fmt.Sprintf("processed %d lines of file %s",
						count, fname))
		}
	}
	removeProgress(fname)
	fmt.Println(fmt.Sprintf("finished processing %d lines of file %s in %s",
				count, fname, util.ElapsedTime(start, time.Now())))
	return count
}


func processFile(fname string) string {
	dname := filepath.Dir(fname)
	bname := filepath.Base(fname)
	return filepath.Join(dname, fmt.Sprintf(".%s.progress", bname))
}


func getProcessedLines(fname string) int {
	var err error
	fname = processFile(fname)
	if _, err = os.Stat(fname); os.IsNotExist(err) {
		return 0
	}
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return 0
	}
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return 0
	}
	return i

}


func storeProgress(fname string, lines int) bool {
	fname = processFile(fname)
	err := ioutil.WriteFile(fname, []byte(fmt.Sprintf("%d", lines)), 0644)
	if err != nil {
		return false
	}
	return true
}


func removeProgress(fname string) bool {
	var err error
	fname = processFile(fname)
	if _, err = os.Stat(fname); os.IsNotExist(err) {
		return true
	}
	err = os.Remove(fname)
	if err != nil {
		return false
	}
	return true
}


func main() {
	config := cfg.ReadArgs()
	if len(config.Files) == 0 {
		fmt.Println("you must specify at least one file to process")
		os.Exit(1)
	}
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
