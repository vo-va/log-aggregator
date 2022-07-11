package pipeline

import (
	"time"

	"log-aggregator/cursor"
	"log-aggregator/destinations"
	"log-aggregator/sources"
	"log-aggregator/transform"
	"log-aggregator/types"

	"github.com/cenkalti/backoff"
)

type Pipeline struct {
	progress chan types.Cursor
	input    chan *types.Record
	output   chan *types.Record
	conf     Config
}

type Config struct {
	MaxBuffer    int
	Cursor       cursor.DB
	Input        sources.Source
	Destination  destinations.Destination
	Transformers []transform.Transformer
}

func New(conf Config) (*Pipeline, error) {
	input := make(chan *types.Record, conf.MaxBuffer)
	output := make(chan *types.Record, 20)
	progress := make(chan types.Cursor, 5)

	return &Pipeline{
		input:    input,
		output:   output,
		progress: progress,
		conf:     conf,
	}, nil
}

func (p *Pipeline) Start() {
	p.conf.Input.Start(p.input)
	p.conf.Destination.Start(p.output, p.progress)
	go p.transform()
	go p.syncCursor()
}

func (p *Pipeline) Stop(timeout time.Duration) {
	p.conf.Input.Stop()
	time.Sleep(timeout)
}

func (p *Pipeline) transform() {
	for {
		record, open := <-p.input
		if !open {
			return
		}

		for _, transformer := range p.conf.Transformers {
			record, _ = transformer(record)
		}
		p.output <- record
	}
}

func (p *Pipeline) syncCursor() {
	for {
		cursor, open := <-p.progress
		if !open {
			return
		}
		strategy := backoff.NewExponentialBackOff()
		strategy.MaxElapsedTime = time.Second * 15
		err := backoff.Retry(func() error {
			return p.conf.Cursor.Set(cursor)
		}, strategy)
		if err != nil {
			panic(err)
		}
	}
}
