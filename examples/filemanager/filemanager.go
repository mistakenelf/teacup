package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/code"
	"github.com/knipferrc/teacup/filetree"
	"github.com/knipferrc/teacup/help"
	"github.com/knipferrc/teacup/image"
	"github.com/knipferrc/teacup/markdown"
	"github.com/knipferrc/teacup/pdf"
	"github.com/knipferrc/teacup/statusbar"
)

type sessionState int

const (
	idleState sessionState = iota
	showCodeState
	showImageState
	showMarkdownState
	showPdfState
)

// Bubble represents the properties of the UI.
type Bubble struct {
	filetree  filetree.Bubble
	help      help.Bubble
	code      code.Bubble
	image     image.Bubble
	markdown  markdown.Bubble
	pdf       pdf.Bubble
	statusbar statusbar.Bubble
	state     sessionState
	activeBox int
}

// New creates a new instance of the UI.
func New() Bubble {
	filetreeModel := filetree.New(
		true,
		false,
		lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "63", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
	)
	codeModel := code.New(false, false, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	codeModel.SetSyntaxTheme("pygments")
	imageModel := image.New(false, false, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	markdownModel := markdown.New(false, false, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	pdfModel := pdf.New(false, false, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	helpModel := help.New(
		false,
		false,
		"Help",
		help.TitleColor{
			Background: lipgloss.AdaptiveColor{Light: "62", Dark: "62"},
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
		},
		lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"},
		[]help.Entry{
			{Key: "ctrl+c", Description: "Exit FM"},
			{Key: "j/up", Description: "Move up"},
			{Key: "k/down", Description: "Move down"},
			{Key: "h/left", Description: "Go back a directory"},
			{Key: "l/right", Description: "Read file or enter directory"},
			{Key: "p", Description: "Preview directory"},
			{Key: "G", Description: "Jump to bottom"},
			{Key: "~", Description: "Go to home directory"},
			{Key: ".", Description: "Toggle hidden files"},
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
			{Key: "tab", Description: "Toggle between boxes"},
		},
	)

	return Bubble{
		filetree:  filetreeModel,
		help:      helpModel,
		code:      codeModel,
		image:     imageModel,
		markdown:  markdownModel,
		pdf:       pdfModel,
		statusbar: statusbar.Bubble{},
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

	b.filetree, cmd = b.filetree.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		resizeImgCmd := b.image.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		markdownCmd := b.markdown.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		b.filetree.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		b.help.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		b.code.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		b.pdf.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		b.statusbar.SetSize(msg.Width)

		cmds = append(cmds, resizeImgCmd, markdownCmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return b, tea.Quit
		case "q":
			if !b.filetree.IsFiltering() {
				return b, tea.Quit
			}
		case " ":
			selectedFile := b.filetree.GetSelectedItem()
			if !selectedFile.IsDirectory() {
				if selectedFile.FileExtension() == ".png" || selectedFile.FileExtension() == ".jpg" {
					b.state = showImageState
					readFileCmd := b.image.SetFileName(selectedFile.FileName())
					cmds = append(cmds, readFileCmd)
				} else if selectedFile.FileExtension() == ".md" {
					b.state = showMarkdownState
					markdownCmd := b.markdown.SetFileName(selectedFile.FileName())
					cmds = append(cmds, markdownCmd)
				} else if selectedFile.FileExtension() == ".pdf" {
					b.state = showPdfState
					pdfCmd := b.pdf.SetFileName(selectedFile.FileName())
					cmds = append(cmds, pdfCmd)
				} else {
					b.state = showCodeState
					readFileCmd := b.code.SetFileName(selectedFile.FileName())
					cmds = append(cmds, readFileCmd)
				}
			}
		case "tab":
			b.activeBox = (b.activeBox + 1) % 2
			if b.activeBox == 0 {
				b.filetree.SetIsActive(true)
				b.code.SetIsActive(false)
				b.markdown.SetIsActive(false)
				b.image.SetIsActive(false)
				b.pdf.SetIsActive(false)
				b.help.SetIsActive(false)
				b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
				b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				b.markdown.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
			} else {
				switch b.state {
				case idleState:
					b.filetree.SetIsActive(false)
					b.code.SetIsActive(false)
					b.markdown.SetIsActive(false)
					b.image.SetIsActive(false)
					b.pdf.SetIsActive(false)
					b.help.SetIsActive(true)
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.markdown.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				case showCodeState:
					b.filetree.SetIsActive(false)
					b.code.SetIsActive(true)
					b.markdown.SetIsActive(false)
					b.image.SetIsActive(false)
					b.pdf.SetIsActive(false)
					b.help.SetIsActive(false)
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.markdown.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				case showImageState:
					b.filetree.SetIsActive(false)
					b.code.SetIsActive(false)
					b.markdown.SetIsActive(false)
					b.image.SetIsActive(true)
					b.pdf.SetIsActive(false)
					b.help.SetIsActive(false)
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.markdown.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
				case showMarkdownState:
					b.filetree.SetIsActive(false)
					b.code.SetIsActive(false)
					b.markdown.SetIsActive(true)
					b.image.SetIsActive(false)
					b.pdf.SetIsActive(false)
					b.help.SetIsActive(false)
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.markdown.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
				}
			}
		}
	}

	b.statusbar.SetContent(
		b.filetree.GetSelectedItem().ShortName(),
		b.filetree.GetSelectedItem().CurrentDirectory(),
		fmt.Sprintf("%d/%d", b.filetree.Cursor(), b.filetree.TotalItems()),
		"FM",
	)

	b.code, cmd = b.code.Update(msg)
	cmds = append(cmds, cmd)

	b.markdown, cmd = b.markdown.Update(msg)
	cmds = append(cmds, cmd)

	b.image, cmd = b.image.Update(msg)
	cmds = append(cmds, cmd)

	b.pdf, cmd = b.pdf.Update(msg)
	cmds = append(cmds, cmd)

	b.help, cmd = b.help.Update(msg)
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
	case showPdfState:
		rightBox = b.pdf.View()
	case showMarkdownState:
		rightBox = b.markdown.View()
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox),
		b.statusbar.View(),
	)
}

func main() {
	b := New()
	p := tea.NewProgram(b, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
