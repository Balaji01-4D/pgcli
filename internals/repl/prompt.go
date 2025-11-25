package repl

import (
	"context"
	"fmt"
	"os"
	"pgcli/internals/database"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

type Model struct {
	ctx       context.Context
	exec      database.Executor
	textInput textinput.Model
	err       error
}

func InitialModel(ctx context.Context, executor database.Executor) Model {
	ti := textinput.New()
	ti.Placeholder = "Type something..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return Model{
		ctx:       ctx,
		exec:      executor,
		textInput: ti,
		err:       nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			input := m.textInput.Value()
			if input == "" {
				return m, nil
			}

			// Format the output to be printed.
			result, err := m.exec.Execute(m.ctx, input)
			if err != nil {
				m.err = err
				m.textInput.Reset()
				err := fmt.Sprintf("> %s\n%s", input, errorStyle.Render(err.Error()))
				return m, tea.Println(err)
			}

			var out string

			switch r := result.(type) {
			case *database.QueryResult:
				out = r.RenderTable()
			case *database.ExecResult:
				out = r.Render()
			}
			output := fmt.Sprintf("> %s\n%s", input, out)

			// tea.Println prints the string above our view.
			printCmd := tea.Println(output)

			m.textInput.Reset()

			// Return the command to be executed by the Bubble Tea runtime.
			return m, printCmd

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	// The view now only consists of the text input and a quit message.
	// The history is printed to the terminal, not managed by the view.
	return m.textInput.View()
}

func StartREPL(ctx context.Context, executor database.Executor) {
	p := tea.NewProgram(InitialModel(ctx, executor))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting REPL:", err)
		os.Exit(1)
	}
}
