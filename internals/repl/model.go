package repl

import (
	"pgcli/internals/completer"

	"github.com/elk-language/go-prompt"
)


type Repl struct {
	db string
}

func (r *Repl) GetPrefix() string {
	return r.db + "> "
}

func (r *Repl) GetLine() string {
	text := prompt.Input(
		prompt.WithPrefix(r.GetPrefix()),
		prompt.WithCompleter(completer.SQLKeywordCompleter),
	)
	return text
}

func NewModel(db string) *Repl {
	return &Repl{db: db}
}



