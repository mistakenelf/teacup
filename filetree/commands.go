package filetree

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistakenelf/teacup/filesystem"
)

type getDirectoryListingMsg []DirectoryItem
type errorMsg error

// getDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func getDirectoryListingCmd(directoryName string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		var directoryItems []DirectoryItem

		if directoryName == filesystem.HomeDirectory {
			directoryName, err = filesystem.GetHomeDirectory()
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

		files, err := filesystem.GetDirectoryListing(directoryName, showHidden)
		if err != nil {
			return errorMsg(err)
		}

		err = os.Chdir(directoryName)
		if err != nil {
			return errorMsg(err)
		}

		workingDirectory, err := filesystem.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err)
		}

		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}

			status := fmt.Sprintf("%s %s %s",
				fileInfo.ModTime().Format("2006-01-02 15:04:05"),
				fileInfo.Mode().String(),
				ConvertBytesToSizeString(fileInfo.Size()))

			directoryItems = append(directoryItems, DirectoryItem{
				name:             file.Name(),
				details:          status,
				path:             filepath.Join(workingDirectory, file.Name()),
				extension:        filepath.Ext(fileInfo.Name()),
				isDirectory:      fileInfo.IsDir(),
				currentDirectory: workingDirectory,
			})
		}

		return getDirectoryListingMsg(directoryItems)
	}
}
