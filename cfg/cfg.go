package cfg

import (
	"flag"
)

type Config struct {
	MappingFile string
	Index string
	Type string
	Source int
	PartitionLevels int
	PartitionChars string
	ElasticsearchURL string
	DocumentSchemaFile string
	WriterWorkers int
	LineWorkers int
	BufferedLines int
	BulkActions int
	BulkSize int
	Files []string
}


func ReadArgs() (*Config) {
	c := Config{}
	partitionChars := "abcdefghijklmnopqrstuvwxyz"
	partitionChars += "0123456789"
	flag.StringVar(&c.Index, "index", "pwd_",
			"name (or prefix) of the index")
	flag.StringVar(&c.Type, "type", "account",
			"type of the documents")
	flag.IntVar(&c.Source, "source", 1,
			"a number representing the source of the documents")
	flag.IntVar(&c.PartitionLevels, "partition-levels", 1,
			"levels of [a-zA-Z0-9] (plus 'misc') partitions to create")
	flag.StringVar(&c.PartitionChars, "partition-chars", partitionChars,
			"chars to use to generate partition combinations")
	flag.StringVar(&c.ElasticsearchURL, "elasticsearch-url", "http://127.0.0.1:9200",
			"URL of the Elasticsearch server or cluster")
	// TODO: read the mapping from a JSON file.
	//       (insert a placeholder for the document type)
	flag.StringVar(&c.MappingFile, "mapping-file", "",
			"file containing the mapping settings to use")
	// TODO: read the document schema from a JSON file.
	flag.StringVar(&c.DocumentSchemaFile, "document-schema-file", "",
			"file containing the document schema")
	flag.IntVar(&c.WriterWorkers, "writer-workers", 5,
			"number of writer workers to spawn")
	flag.IntVar(&c.LineWorkers, "line-workers", 10,
			"number of line workers to spawn")
	flag.IntVar(&c.BufferedLines, "buffered-lines", 10000,
			"number of buffered lines")
	flag.IntVar(&c.BulkActions, "bulk-actions", 100000,
			"number of insert of the bulk request before a flush")
	flag.IntVar(&c.BulkSize, "bulk-size", 5,
			"megabytes of the bulk request before a flush")
	flag.Parse()
	c.Files = flag.Args()
	return &c
}

