package filetree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knipferrc/teacup/dirfs"
)

type getDirectoryListingMsg []list.Item
type errorMsg error
type copyToClipboardMsg string
type editorFinishedMsg struct{ err error }

// getDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func getDirectoryListingCmd(directoryName string, showHidden, showIcons bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		var items []list.Item

		if directoryName == dirfs.HomeDirectory {
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

		workingDirectory, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err)
		}

		items = append(items, Item{
			title:            dirfs.PreviousDirectory,
			desc:             "",
			shortName:        dirfs.PreviousDirectory,
			fileName:         filepath.Join(workingDirectory, dirfs.PreviousDirectory),
			extension:        "",
			isDirectory:      directoryInfo.IsDir(),
			currentDirectory: workingDirectory,
			fileInfo:         nil,
			showIcons:        false,
		})

		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}

			status := fmt.Sprintf("%s %s %s",
				fileInfo.ModTime().Format("2006-01-02 15:04:05"),
				fileInfo.Mode().String(),
				ConvertBytesToSizeString(fileInfo.Size()))

			items = append(items, Item{
				title:            file.Name(),
				desc:             status,
				shortName:        file.Name(),
				fileName:         filepath.Join(workingDirectory, file.Name()),
				extension:        filepath.Ext(fileInfo.Name()),
				isDirectory:      fileInfo.IsDir(),
				currentDirectory: workingDirectory,
				fileInfo:         fileInfo,
				showIcons:        showIcons,
			})
		}

		return getDirectoryListingMsg(items)
	}
}

// moveItemCmd moves files to the current directory.
func moveItemCmd(path, name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err)
		}

		if err := dirfs.MoveDirectoryItem(path, fmt.Sprintf("%s/%s", workingDir, name)); err != nil {
			return errorMsg(err)
		}

		return nil
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
		fileInfo, err := os.Lstat(name)
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

		return copyToClipboardMsg(fmt.Sprintf(
			"%s %s %s",
			"Successfully copied", name, "to clipboard",
		))
	}
}

// renameItemCmd renames a file or directory based on the name and value provided.
func renameItemCmd(name, value string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.RenameDirectoryItem(name, value); err != nil {
			return errorMsg(err)
		}

		return nil
	}
}

// writeSelectionPathCmd writes content to the file specified.
func writeSelectionPathCmd(selectionPath, filePath string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.WriteToFile(selectionPath, filePath); err != nil {
			return errorMsg(err)
		}

		return nil
	}
}

// openInEditor opens the file in the editor specified and default to vim if not set.
func openInEditor(fileName string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	c := exec.Command(editor, fileName) //nolint:gosec

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}
