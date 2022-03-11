// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/dirfs"
)

// SetSize sets the size of the filetree.
func (b *Bubble) SetSize(width, height int) {
	v, h := listStyle.GetFrameSize()

	b.list.Styles.StatusBar.Width(width - h)
	b.list.SetSize(width-h, height-v-lipgloss.Height(b.input.View())-inputStyle.GetVerticalPadding())
}

// SetBorderColor sets the color of the border.
func (b *Bubble) SetBorderColor(color lipgloss.AdaptiveColor) {
	listStyle = listStyle.Copy().BorderForeground(color)
}

// GetSelectedItem returns the currently selected item in the tree.
func (b Bubble) GetSelectedItem() item {
	selectedDir, ok := b.list.SelectedItem().(item)
	if ok {
		return selectedDir
	}

	return item{}
}

// Update handles updating the filetree.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case getDirectoryListingMsg:
		cmd = b.list.SetItems(msg)
		cmds = append(cmds, cmd)
	case errorMsg:
		return b, b.list.NewStatusMessage(statusMessageErrorStyle(msg.Error()))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, openDirectoryKey):
			if !b.input.Focused() {
				selectedDir := b.GetSelectedItem()
				cmds = append(cmds, getDirectoryListingCmd(selectedDir.fileName, b.showHidden))
			}
		case key.Matches(msg, copyItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully copied file"),
				)

				cmds = append(cmds, tea.Sequentially(
					copyItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)
			}
		case key.Matches(msg, zipItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully zipped item"),
				)

				cmds = append(cmds, tea.Sequentially(
					zipItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)
			}
		case key.Matches(msg, unzipItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully unzipped item"),
				)

				cmds = append(cmds, tea.Sequentially(
					unzipItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)
			}
		case key.Matches(msg, createFileKey):
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Enter name of new file"
				b.state = createFileState

				return b, nil
			}
		case key.Matches(msg, createDirectoryKey):
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Enter name of new directory"
				b.state = createDirectoryState

				return b, nil
			}
		case key.Matches(msg, deleteItemKey):
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Are you sure you want to delete (y/n)?"
				b.state = deleteItemState

				return b, nil
			}
		case key.Matches(msg, toggleHiddenKey):
			if !b.input.Focused() {
				b.showHidden = !b.showHidden
				cmds = append(cmds, getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden))
			}
		case key.Matches(msg, submitInputKey):
			switch b.state {
			case createFileState:
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully created file"),
				)

				cmds = append(cmds, tea.Sequentially(
					createFileCmd(b.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)

				b.input.Blur()
				b.input.Reset()
			case createDirectoryState:
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully created directory"),
				)

				cmds = append(cmds, statusCmd)
				cmds = append(cmds, tea.Sequentially(
					createDirectoryCmd(b.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))

				b.input.Blur()
				b.input.Reset()
			case deleteItemState:
				if strings.ToLower(b.input.Value()) == "y" {
					selectedDir := b.GetSelectedItem()

					statusCmd := b.list.NewStatusMessage(
						statusMessageInfoStyle("Successfully deleted item"),
					)

					cmds = append(cmds, statusCmd)
					cmds = append(cmds, tea.Sequentially(
						deleteItemCmd(selectedDir.fileName),
						getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
					))
				}

				b.input.Blur()
				b.input.Reset()
			}
		}
	}

	b.list, cmd = b.list.Update(msg)
	cmds = append(cmds, cmd)

	b.input, cmd = b.input.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}
