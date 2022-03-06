package filetree

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/dirfs"
)

type getDirectorylistingMsg []list.Item
type errorMsg string

var (
	listStyle              = lipgloss.NewStyle().Margin(1)
	statusMessageInfoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	statusMessageErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}).
				Render
)

// getDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func getDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := dirfs.GetDirectoryListing(name, true)
		if err != nil {
			return errorMsg(err.Error())
		}

		err = os.Chdir(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		var items []list.Item
		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}

			status := fmt.Sprintf("%s %s",
				fileInfo.ModTime().Format("2006-01-02 15:04:05"),
				fileInfo.Mode().String())

			items = append(items, item{
				title: file.Name(),
				desc:  status,
			})
		}

		return getDirectorylistingMsg(items)
	}
}

// Bubble represents the properties of a filetree.
type Bubble struct {
	list list.Model
}

// item represents a list item.
type item struct {
	title, desc string
}

// Title returns the title of the list item.
func (i item) Title() string {
	return i.title
}

// Description returns the description of the list item.
func (i item) Description() string { return i.desc }

// FilterValue returns the current filter value.
func (i item) FilterValue() string { return i.title }

func New() Bubble {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Filetree"

	return Bubble{
		list: l,
	}
}

func (b Bubble) Init() tea.Cmd {
	return getDirectoryListingCmd(dirfs.CurrentDirectory)
}

// Update handles updating the filetree.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case getDirectorylistingMsg:
		b.list.SetItems(msg)

		return b, nil
	case tea.WindowSizeMsg:
		v, h := listStyle.GetFrameSize()
		b.list.SetSize(msg.Width-h, msg.Height-v)
	}

	b.list, cmd = b.list.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

// View returns a string representation of a filetree.
func (b Bubble) View() string {
	return listStyle.Render(b.list.View())
}
