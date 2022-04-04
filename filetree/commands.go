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
	"github.com/knipferrc/teacup/dirfs"
)

type getDirectoryListingMsg []list.Item
type errorMsg error
type copyToClipboardMsg string

// getDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func getDirectoryListingCmd(name string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		var items []list.Item

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

		workingDirectory, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err)
		}

		items = append(items, item{
			title:            dirfs.PreviousDirectory,
			desc:             "",
			shortName:        dirfs.PreviousDirectory,
			fileName:         filepath.Join(workingDirectory, dirfs.PreviousDirectory),
			extension:        "",
			isDirectory:      directoryInfo.IsDir(),
			currentDirectory: workingDirectory,
			fileInfo:         nil,
		})

		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}

			if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
				symlinkFile, err := os.Readlink(fileInfo.Name())
				if err != nil {
					return errorMsg(err)
				}

				symlinkFileInfo, err := os.Stat(symlinkFile)
				if err != nil {
					return errorMsg(err)
				}

				status := fmt.Sprintf("%s %s %s",
					symlinkFileInfo.ModTime().Format("2006-01-02 15:04:05"),
					symlinkFileInfo.Mode().String(),
					ConvertBytesToSizeString(symlinkFileInfo.Size()))

				items = append(items, item{
					title:            fileInfo.Name(),
					desc:             status,
					shortName:        fileInfo.Name(),
					fileName:         filepath.Join(workingDirectory, symlinkFileInfo.Name()),
					extension:        filepath.Ext(symlinkFileInfo.Name()),
					isDirectory:      symlinkFileInfo.IsDir(),
					currentDirectory: workingDirectory,
					fileInfo:         fileInfo,
				})
			} else {
				status := fmt.Sprintf("%s %s %s",
					fileInfo.ModTime().Format("2006-01-02 15:04:05"),
					fileInfo.Mode().String(),
					ConvertBytesToSizeString(fileInfo.Size()))

				items = append(items, item{
					title:            file.Name(),
					desc:             status,
					shortName:        file.Name(),
					fileName:         filepath.Join(workingDirectory, file.Name()),
					extension:        filepath.Ext(fileInfo.Name()),
					isDirectory:      fileInfo.IsDir(),
					currentDirectory: workingDirectory,
					fileInfo:         fileInfo,
				})
			}
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

// handleErrorCmd returns an error message to the UI.
func handleErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg(err)
	}
}

// redrawCmd redraws the UI.
func (b Bubble) redrawCmd() tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  b.width,
			Height: b.height,
		}
	}
}
