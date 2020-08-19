package elktail

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"
)

// Tail @todo
type Tail interface {
	Run(context.Context, io.Writer) error
	Stop() error // is this necessary? i don't think it is anymore
}

// Config @todo
type Config struct {
	After        *time.Time
	PageSize     int
	IndexPattern string
	ListOnly     bool
	Query        string
	URL          url.URL
}

// NewTail @todo
func NewTail(c Config) Tail {
	return &tail{}

}

type tail struct {
	out    chan []byte
	config Config
	offset *time.Time
}

func (t *tail) Run(ctx context.Context, out io.Writer) error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-ticker.C:
			_, err := fmt.Fprintln(out, "{\"tick tock\":\"", t, "\"}")
			if err != nil {
				return err
			}
		}
	}
}

func (t *tail) Stop() error {
	return nil
}
