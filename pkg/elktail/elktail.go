package elktail

import (
	"context"
	"encoding/json"
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
	return &tail{config: c}

}

type tail struct {
	out    chan []byte
	config Config
	offset *time.Time
}

type status struct {
	Timestamp time.Time
	Config
}

func (t *tail) Run(ctx context.Context, out io.Writer) error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case i := <-ticker.C:
			s := status{Config: t.config, Timestamp: i}
			b, err := json.Marshal(s)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(out, string(b))
			if err != nil {
				return err
			}
		}
	}
}

func (t *tail) Stop() error {
	return nil
}
