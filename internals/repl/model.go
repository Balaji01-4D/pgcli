package repl

import (
	"fmt"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/knz/bubbline"
)


type Model struct {
	*bubbline.Editor
	CurrentDB string
}

func NewModel(db string) *Model {
	lineedit := bubbline.New()
	lineedit.CursorMode = cursor.CursorBlink
	lineedit.Prompt = fmt.Sprintf("%s> ", db)

	return &Model{
		Editor:   lineedit,
		CurrentDB: db,
	}
}



