package repl

import (
	"io"
	"os"
	"time"

	"github.com/elk-language/go-prompt"
	"github.com/fatih/color"
)

type Repl struct {
	db string
}

func NewModel(db string) *Repl {
	return &Repl{db: db}
}

func (r *Repl) GetPrefix() string {
	return r.db + "> "
}

func (r *Repl) Read() string {
	text := prompt.Input(
		prompt.WithPrefix(r.GetPrefix()),
	)
	return text
}

func (r *Repl) PrintError(err error) {
	c := color.New(color.FgRed)
	c.Fprintf(os.Stderr, "%v\n", err)
}

func (r *Repl) PrintTime(time time.Duration) {
	c := color.New(color.FgCyan)
	c.Fprintf(os.Stderr, "Time: %.3fs\n", time.Seconds())
}

func (r *Repl) Print(output string) {
	EchoViaPager(func(w io.Writer) error {
		io.WriteString(w, output)
		return nil
	})
}
