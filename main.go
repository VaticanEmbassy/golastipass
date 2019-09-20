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

type Processor struct {
	config *cfg.Config
	writer *elastic.Writer
	linesCh chan string
	wg sync.WaitGroup
}


func (p* Processor) ProcessLine() {
	for {
		line, ok := <- p.linesCh
		if ok {
			d, err := doc.ParseLine(line, p.config.Type, p.config.Source,
			p.config.PartitionLevels)
			if err == nil {
				p.writer.Add(&d)
			}
		} else {
			p.wg.Done()
			return
		}
	}
}


func (p* Processor) ProcessFile(fname string) int {
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
		linesToSkip = p.getProgress(fname)
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
	canRestore := p.createProgress(fname)
	if !canRestore {
		fmt.Printf("unable to write %s file; operation can't be stopped and restored",
			p.progressFilePath(fname))
	}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		count += 1
		p.linesCh <- line
		if (saveProgress && count % 10000 == 0) {
			p.storeProgress(fname, count)
		}
		if (count % 100000 == 0) {
			fmt.Println(fmt.Sprintf("processed %d lines of file %s",
						count, fname))
		}
	}
	p.removeProgress(fname)
	fmt.Println(fmt.Sprintf("finished processing %d lines of file %s in %s",
				count, fname, util.ElapsedTime(start, time.Now())))
	return count
}


func (p* Processor) progressFilePath(fname string) string {
	dname := filepath.Dir(fname)
	bname := filepath.Base(fname)
	return filepath.Join(dname, fmt.Sprintf(".%s.progress", bname))
}


func (p* Processor) getProgress(fname string) int {
	var err error
	fname = p.progressFilePath(fname)
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


func (p* Processor) storeProgress(fname string, lines int) bool {
	fname = p.progressFilePath(fname)
	err := ioutil.WriteFile(fname, []byte(fmt.Sprintf("%d", lines)), 0644)
	if err != nil {
		return false
	}
	return true
}


func (p* Processor) createProgress(fname string) bool {
	f, err := os.OpenFile(p.progressFilePath(fname),
				os.O_APPEND |os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	defer f.Close()
	if _, err := f.WriteString(""); err != nil {
		return false
	}
	return true
}


func (p* Processor) removeProgress(fname string) bool {
	var err error
	fname = p.progressFilePath(fname)
	if _, err = os.Stat(fname); os.IsNotExist(err) {
		return true
	}
	err = os.Remove(fname)
	if err != nil {
		return false
	}
	return true
}


func (p* Processor) CreateIndices() {
	start := time.Now()
	fmt.Print("creating indices, please wait... ")
	indicesNr := p.writer.CreateIndices()
	fmt.Printf("created %d indices in %s.\n", indicesNr,
		util.ElapsedTime(start, time.Now()))
}


func (p* Processor) Run() {
	p.writer = elastic.NewWriter(p.config)
	p.CreateIndices()
	start := time.Now()

	p.linesCh = make(chan string, p.config.BufferedLines)
	for w := 1; w <= p.config.LineWorkers; w++ {
		p.wg.Add(1)
		go p.ProcessLine()
	}

	count := 0
	for _, fname := range p.config.Files {
		count += 1
		p.ProcessFile(fname)
	}

	close(p.linesCh)
	p.wg.Wait()
	p.writer.Close()
	fmt.Printf("processed %d files in %v\n",
		count, util.ElapsedTime(start, time.Now()))
}


func NewProcessor(config *cfg.Config) *Processor {
	p := new(Processor)
	p.config = config
	return p
}


func main() {
	config := cfg.ReadArgs()
	if len(config.Files) == 0 {
		fmt.Println("you must specify at least one file to process")
		os.Exit(1)
	}
	p := NewProcessor(config)
	p.Run()
}
