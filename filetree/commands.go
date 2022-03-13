// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/dirfs"
	"github.com/knipferrc/teacup/formatter"
	"github.com/knipferrc/teacup/icons"
)

type getDirectoryListingMsg []list.Item
type errorMsg error
type copyToClipboardMsg string

// getDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func getDirectoryListingCmd(name string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		directoryName := name

		if name == dirfs.HomeDirectory {
			directoryName, err = dirfs.GetHomeDirectory()
			if err != nil {
				return errorMsg(err)
			}
		}

		directoryInfo, err := os.Stat(directoryName)
		if err != nil {
			return errorMsg(err)
		}

		if !directoryInfo.IsDir() {
			return nil
		}

		files, err := dirfs.GetDirectoryListing(directoryName, showHidden)
		if err != nil {
			return errorMsg(err)
		}

		err = os.Chdir(directoryName)
		if err != nil {
			return errorMsg(err)
		}

		var items []list.Item

		workingDirectory, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err)
		}

		items = append(items, Item{
			ItemTitle: dirfs.PreviousDirectory,
			Desc:      "",
			FileName:  filepath.Join(workingDirectory, dirfs.PreviousDirectory),
		})

		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}

			icon, color := icons.GetIcon(fileInfo.Name(), filepath.Ext(fileInfo.Name()), icons.GetIndicator(fileInfo.Mode()))
			fileIcon := lipgloss.NewStyle().Width(fileIconWidth).Render(fmt.Sprintf("%s%s\033[0m ", color, icon))

			status := fmt.Sprintf("%s %s %s",
				fileInfo.ModTime().Format("2006-01-02 15:04:05"),
				fileInfo.Mode().String(),
				formatter.ConvertBytesToSizeString(fileInfo.Size()))

			items = append(items, Item{
				ItemTitle: lipgloss.JoinHorizontal(lipgloss.Top, fileIcon, file.Name()),
				Desc:      status,
				FileName:  filepath.Join(workingDirectory, file.Name()),
			})
		}

		return getDirectoryListingMsg(items)
	}
}

// createFileCmd creates a file based on the name provided.
func createFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CreateFile(name); err != nil {
			return errorMsg(err)
		}

		return nil
	}
}

// createDirectoryCmd creates a directory based on the name provided.
func createDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CreateDirectory(name); err != nil {
			return errorMsg(err)
		}

		return nil
	}
}

// deleteDirectoryCmd deletes a directory based on the name provided.
func deleteItemCmd(name string) tea.Cmd {
	return func() tea.Msg {
		fileInfo, err := os.Stat(name)
		if err != nil {
			return errorMsg(err)
		}

		if fileInfo.IsDir() {
			if err := dirfs.DeleteDirectory(name); err != nil {
				return errorMsg(err)
			}
		} else {
			if err := dirfs.DeleteFile(name); err != nil {
				return errorMsg(err)
			}
		}

		return nil
	}
}

// zipItemCmd zips a directory based on the name provided.
func zipItemCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.Zip(name); err != nil {
			return errorMsg(err)
		}

		return nil
	}
}

// unzipItemCmd unzips a directory based on the name provided.
func unzipItemCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.Unzip(name); err != nil {
			return errorMsg(err)
		}

		return nil
	}
}

// copyItemCmd copies a file or directory given a name.
func copyItemCmd(name string) tea.Cmd {
	return func() tea.Msg {
		fileInfo, err := os.Stat(name)
		if err != nil {
			return errorMsg(err)
		}

		if fileInfo.IsDir() {
			if err := dirfs.CopyDirectory(name); err != nil {
				return errorMsg(err)
			}
		} else {
			if err := dirfs.CopyFile(name); err != nil {
				return errorMsg(err)
			}
		}

		return nil
	}
}

// copyToClipboardCmd copies the provided string to the clipboard.
func copyToClipboardCmd(name string) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(name)
		if err != nil {
			return errorMsg(err)
		}

		return copyToClipboardMsg(fmt.Sprintf("%s %s %s", "Successfully copied", name, "to clipboard"))
	}
}
