package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/code"
	"github.com/knipferrc/teacup/filetree"
	"github.com/knipferrc/teacup/help"
	"github.com/knipferrc/teacup/image"
)

type sessionState int

const (
	idleState sessionState = iota
	showCodeState
	showImageState
)

// Bubble represents the properties of the UI.
type Bubble struct {
	filetree  filetree.Bubble
	help      help.Bubble
	code      code.Bubble
	image     image.Bubble
	state     sessionState
	activeBox int
}

// New creates a new instance of the UI.
func New() Bubble {
	filetreeModel := filetree.New(false, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	codeModel := code.New(false, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	imageModel := image.New(false, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})

	filetreeModel.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})

	helpModel := help.New(
		lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"},
		"Help",
		[]help.Entry{
			{Key: "ctrl+c", Description: "Exit FM"},
			{Key: "j/up", Description: "Move up"},
			{Key: "k/down", Description: "Move down"},
			{Key: "h/left", Description: "Go back a directory"},
			{Key: "l/right", Description: "Read file or enter directory"},
			{Key: "p", Description: "Preview directory"},
			{Key: "gg", Description: "Go to top of filetree or box"},
			{Key: "G", Description: "Go to bottom of filetree or box"},
			{Key: "~", Description: "Go to home directory"},
			{Key: "/", Description: "Go to root directory"},
			{Key: ".", Description: "Toggle hidden files"},
			{Key: "S", Description: "Only show directories"},
			{Key: "s", Description: "Only show files"},
			{Key: "y", Description: "Copy file path to clipboard"},
			{Key: "Z", Description: "Zip currently selected tree item"},
			{Key: "U", Description: "Unzip currently selected tree item"},
			{Key: "n", Description: "Create new file"},
			{Key: "N", Description: "Create new directory"},
			{Key: "ctrl+d", Description: "Delete currently selected tree item"},
			{Key: "M", Description: "Move currently selected tree item"},
			{Key: "enter", Description: "Process command"},
			{Key: "E", Description: "Edit currently selected tree item"},
			{Key: "C", Description: "Copy currently selected tree item"},
			{Key: "esc", Description: "Reset FM to initial state"},
			{Key: "O", Description: "Show logs if debugging enabled"},
			{Key: "tab", Description: "Toggle between boxes"},
		},
		false,
	)

	return Bubble{
		filetree: filetreeModel,
		help:     helpModel,
		code:     codeModel,
		image:    imageModel,
	}
}

// Init intializes the UI.
func (b Bubble) Init() tea.Cmd {
	return b.filetree.Init()
}

// Update handles all UI interactions.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		resizeImgCmd := b.image.SetSize(msg.Width/2, msg.Height)
		b.filetree.SetSize(msg.Width/2, msg.Height)
		b.help.SetSize(msg.Width/2, msg.Height)
		b.code.SetSize(msg.Width/2, msg.Height)
		cmds = append(cmds, resizeImgCmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return b, tea.Quit
		case " ":
			selectedFile := b.filetree.GetSelectedItem()
			if !selectedFile.IsDirectory() {
				if selectedFile.FileExtension() == ".png" || selectedFile.FileExtension() == ".jpg" {
					b.state = showImageState
					readFileCmd := b.image.SetFileName(selectedFile.FileName())
					cmds = append(cmds, readFileCmd)
				} else {
					b.state = showCodeState
					readFileCmd := b.code.SetFileName(selectedFile.FileName())
					cmds = append(cmds, readFileCmd)
				}
			}
		case "tab":
			b.activeBox = (b.activeBox + 1) % 2
			if b.activeBox == 0 {
				b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
				b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
			} else {
				switch b.state {
				case idleState:
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				case showCodeState:
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				case showImageState:
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
				}
			}
		}
	}

	b.filetree, cmd = b.filetree.Update(msg)
	cmds = append(cmds, cmd)

	b.help, cmd = b.help.Update(msg)
	cmds = append(cmds, cmd)

	b.code, cmd = b.code.Update(msg)
	cmds = append(cmds, cmd)

	b.image, cmd = b.image.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

// View returns a string representation of the UI.
func (b Bubble) View() string {
	leftBox := b.filetree.View()
	rightBox := b.help.View()

	switch b.state {
	case idleState:
		rightBox = b.help.View()
	case showCodeState:
		rightBox = b.code.View()
	case showImageState:
		rightBox = b.image.View()
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
}

func main() {
	b := New()
	p := tea.NewProgram(b, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
