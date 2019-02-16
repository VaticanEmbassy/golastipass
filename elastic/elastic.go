package elastic

import (
	"fmt"
	"log"
	"context"
	"github.com/olivere/elastic"
	"github.com/VaticanEmbassy/golastipass/cfg"
	"github.com/VaticanEmbassy/golastipass/doc"
)


type Writer struct {
	cfg *cfg.Config
	indices map[string]string
	defaultIndex string

	Connection *elastic.Client
	Processor *elastic.BulkProcessor
}


func NewWriter(c *cfg.Config) *Writer {
	w := new(Writer)
	w.cfg = c
	w.Connect()
	w.InitProcessor()
	return w
}


func (w *Writer) Connect() (*elastic.Client) {
	client, err := elastic.NewClient(elastic.SetURL(w.cfg.ElasticsearchURL))
	if err != nil {
		log.Fatal(err)
	}
	w.Connection = client
	return client
}


func (w *Writer) Close() {
	w.Processor.Flush()
	w.Processor.Close()
}


func (w *Writer) CreateIndex(index string) {
	_, err := w.Connection.CreateIndex(index).
			BodyString(cfg.DefaultMapping).Do(context.Background())
	if err != nil {
		e, ok := err.(*elastic.Error)
		if !ok {
			log.Fatal(err)
		}
		if e.Details.Type != "resource_already_exists_exception" {
			log.Fatal(err)
		}
	}
}


func (w *Writer) CreateIndices() int {
	w.indices = make(map[string]string)
	if w.cfg.PartitionLevels == 0 {
		w.defaultIndex = w.cfg.Index
		w.CreateIndex(w.cfg.Index)
		return 1
	}
	w.defaultIndex = w.cfg.Index + "misc"
	w.CreateIndex(w.defaultIndex)
	count := 1
	nc := nextCombination(w.cfg.PartitionLevels, w.cfg.PartitionChars)
	for {
		comb := nc()
		if len(comb) == 0 {
			break
		}
		idx := fmt.Sprintf("%s%s", w.cfg.Index, comb)
		w.indices[comb] = idx
		w.CreateIndex(idx)
		count += 1
	}
	return count
}


func (w *Writer) InitProcessor() (*elastic.BulkProcessor) {
	service := w.Connection.BulkProcessor().Name("golastic-worker-1")
	processor, _ := service.Workers(w.cfg.WriterWorkers).
				After(w.afterCall).
				BulkActions(w.cfg.BulkActions).
				BulkSize(w.cfg.BulkSize << 20).
				Stats(true).
				Do(context.Background())
	w.Processor = processor
	return processor
}

func (w *Writer) afterCall(id int64, requests []elastic.BulkableRequest,
			response *elastic.BulkResponse, err error) {
	if err != nil {
		stats := w.Processor.Stats()
		fmt.Printf("detected insert failure. Global stats: succeeded:%d, failed:%d\n",
			stats.Succeeded, stats.Failed)
	}
}

func (w *Writer) Add(md *doc.MetaDoc) {
	var ok bool
	var index string
	index, ok = w.indices[md.Prefix]
	if !ok {
		index = w.defaultIndex
	}
	d := elastic.NewBulkIndexRequest().Index(index).Type(md.Type).Doc(md.Doc)
	w.Processor.Add(d)
}


/* Stolen from https://stackoverflow.com/a/22741715/253358 */
func nextCombination(n int, c string) func() string {
	r := []rune(c)
	p := make([]rune, n)
	x := make([]int, len(p))
	return func() string {
		p := p[:len(x)]
		for i, xi := range x {
			p[i] = r[xi]
		}
		for i := len(x) - 1; i >= 0; i-- {
			x[i]++
			if x[i] < len(r) {
				break
			}
			x[i] = 0
			if i <= 0 {
				x = x[0:0]
				break
			}
		}
		return string(p)
	}
}
