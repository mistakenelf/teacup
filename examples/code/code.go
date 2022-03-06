package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knipferrc/teacup/code"
)

type Bubble struct {
	code code.Bubble
}

func New() Bubble {
	codeModel := code.New(false)

	return Bubble{
		code: codeModel,
	}
}

func (b Bubble) Init() tea.Cmd {
	return nil
}

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	b.code.SetContent("sourcecode.go")

	b.code, cmd = b.code.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

func (b Bubble) View() string {
	return b.code.View()
}

func main() {
	b := New()
	p := tea.NewProgram(b)

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
