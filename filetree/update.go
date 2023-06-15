package filetree

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knipferrc/teacup/dirfs"
)

const (
	yesKey   = "y"
	enterKey = "enter"
)

// Update handles updating the filetree.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case getDirectoryListingMsg:
		if msg != nil {
			cmd = m.list.SetItems(msg)
			cmds = append(cmds, cmd)
		}
	case copyToClipboardMsg:
		return m, m.list.NewStatusMessage(statusMessageInfoStyle(string(msg)))
	case errorMsg:
		return m, m.list.NewStatusMessage(statusMessageErrorStyle(msg.Error()))
	case tea.KeyMsg:
		if m.IsFiltering() {
			break
		}

		if !m.active {
			return m, nil
		}

		switch m.state {
		case deleteItemState:
			if msg.String() == yesKey {
				selectedItem := m.GetSelectedItem()

				statusCmd := m.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully deleted item"),
				)

				cmds = append(cmds, statusCmd, tea.Sequence(
					deleteItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons),
				))

				m.state = idleState

				return m, tea.Batch(cmds...)
			}
		case moveItemState:
			if msg.String() == enterKey {
				statusCmd := m.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully moved item"),
				)

				cmds = append(cmds, statusCmd, tea.Sequence(
					moveItemCmd(m.itemToMove.path, m.itemToMove.shortName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons),
				))

				m.state = idleState

				return m, tea.Batch(cmds...)
			}
		}

		switch {
		case key.Matches(msg, openDirectoryKey):
			if !m.input.Focused() {
				selectedDir := m.GetSelectedItem()
				cmds = append(cmds, getDirectoryListingCmd(selectedDir.fileName, m.showHidden, m.showIcons))
			}
		case key.Matches(msg, copyItemKey):
			if !m.input.Focused() {
				selectedItem := m.GetSelectedItem()
				statusCmd := m.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully copied file"),
				)

				cmds = append(cmds, statusCmd, tea.Sequence(
					copyItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons),
				))
			}
		case key.Matches(msg, zipItemKey):
			if !m.input.Focused() {
				selectedItem := m.GetSelectedItem()
				statusCmd := m.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully zipped item"),
				)

				cmds = append(cmds, statusCmd, tea.Sequence(
					zipItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons),
				))
			}
		case key.Matches(msg, unzipItemKey):
			if !m.input.Focused() {
				selectedItem := m.GetSelectedItem()
				statusCmd := m.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully unzipped item"),
				)

				cmds = append(cmds, statusCmd, tea.Sequence(
					unzipItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons),
				))
			}
		case key.Matches(msg, createFileKey):
			if !m.input.Focused() {
				m.input.Focus()
				m.input.Placeholder = "Enter name of new file"
				m.state = createFileState

				return m, textinput.Blink
			}
		case key.Matches(msg, createDirectoryKey):
			if !m.input.Focused() {
				m.input.Focus()
				m.input.Placeholder = "Enter name of new directory"
				m.state = createDirectoryState

				return m, textinput.Blink
			}
		case key.Matches(msg, deleteItemKey):
			if !m.input.Focused() {
				m.state = deleteItemState

				return m, nil
			}
		case key.Matches(msg, moveItemKey):
			if !m.input.Focused() {
				selectedItem := m.GetSelectedItem()
				m.state = moveItemState
				m.itemToMove = itemToMove{
					shortName: selectedItem.shortName,
					path:      selectedItem.fileName,
				}

				return m, nil
			}
		case key.Matches(msg, renameItemKey):
			if !m.input.Focused() {
				m.input.Focus()
				m.input.Placeholder = "Enter new name"
				m.state = renameItemState

				return m, textinput.Blink
			}
		case key.Matches(msg, toggleHiddenKey):
			if !m.input.Focused() {
				m.showHidden = !m.showHidden
				cmds = append(cmds, getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons))
			}
		case key.Matches(msg, homeShortcutKey):
			if !m.input.Focused() {
				cmds = append(cmds, getDirectoryListingCmd(dirfs.HomeDirectory, m.showHidden, m.showIcons))
			}
		case key.Matches(msg, rootShortcutKey):
			if !m.input.Focused() {
				cmds = append(cmds, getDirectoryListingCmd(dirfs.RootDirectory, m.showHidden, m.showIcons))
			}
		case key.Matches(msg, copyToClipboardKey):
			if !m.input.Focused() {
				selectedItem := m.GetSelectedItem()
				cmds = append(cmds, copyToClipboardCmd(selectedItem.fileName))
			}
		case key.Matches(msg, escapeKey):
			m.state = idleState

			if m.input.Focused() {
				m.input.Reset()
				m.input.Blur()
			}
		case key.Matches(msg, openInEditorKey):
			if !m.input.Focused() {
				selectedItem := m.GetSelectedItem()

				if m.selectionPath == "" && !selectedItem.IsDirectory() {
					return m, openInEditor(selectedItem.FileName())
				}

				return m, tea.Sequence(
					writeSelectionPathCmd(m.selectionPath, selectedItem.ShortName()),
					tea.Quit,
				)
			}
		case key.Matches(msg, submitInputKey):
			selectedItem := m.GetSelectedItem()

			switch m.state {
			case idleState, deleteItemState, moveItemState:
				return m, nil
			case createFileState:
				statusCmd := m.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully created file"),
				)

				cmds = append(cmds, statusCmd, tea.Sequence(
					createFileCmd(m.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons),
				))
			case createDirectoryState:
				statusCmd := m.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully created directory"),
				)

				cmds = append(cmds, statusCmd, tea.Sequence(
					createDirectoryCmd(m.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons),
				))
			case renameItemState:
				statusCmd := m.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully renamed"),
				)

				cmds = append(cmds, statusCmd, tea.Sequence(
					renameItemCmd(selectedItem.fileName, m.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons),
				))
			}

			m.state = idleState
			m.input.Blur()
			m.input.Reset()
		}
	}

	if m.active {
		switch m.state {
		case idleState, moveItemState:
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		case createFileState, createDirectoryState, renameItemState:
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		case deleteItemState:
			return m, nil
		}
	}

	return m, tea.Batch(cmds...)
}
