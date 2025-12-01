package repl

import "github.com/elk-language/go-prompt"


type Repl struct {
	db string
}

func (r *Repl) GetPrefix() string {
	return r.db
}

func (r *Repl) GetLine() string {
	text := prompt.Input(
		prompt.WithPrefix(r.GetPrefix()),
	)
	return text
}

func NewModel(db string) *Repl {
	return &Repl{db: db}
}



