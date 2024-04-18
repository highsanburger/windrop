package main

import (
	"fmt"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	window := myApp.NewWindow("File Sharing App")

	filePathEntry := widget.NewEntry()
	filePathEntry.SetPlaceHolder("Enter file path")

	codeEntry := widget.NewEntry()
	codeEntry.SetPlaceHolder("Enter code")

	codeRecieve := widget.NewEntry()
	codeRecieve.SetPlaceHolder("Enter code")

	showNotification := func(message string) {
		dialog.ShowInformation("Notification", message, window)
	}

	uploadButton := widget.NewButton("Upload File", func() {
		filePath := filePathEntry.Text
		code := codeEntry.Text

		cmd := exec.Command("croc", "send", "--code", code, filePath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			fmt.Println("error")
			showNotification("Upload failed: " + err.Error())
			return
		}
		fmt.Println(string(output))

		showNotification("File uploaded successfully!")
	})

	downloadButton := widget.NewButton("Download File", func() {
		// Function to handle download
		downloadFunc := func(downloadPath string) {
			code := codeRecieve.Text

			cmd := exec.Command("croc", code)
			cmd.Dir = downloadPath // Set download directory
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
				fmt.Println("error")
				showNotification("Download failed: " + err.Error())
				return
			}
			fmt.Println(string(output))

			showNotification("File downloaded successfully!")
		}

		// Show folder selection dialog
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				downloadPath := uri.Path()
				downloadFunc(downloadPath)
			} else {
				// If no folder selected, download to home directory
				homeDir, err := os.UserHomeDir()
				if err != nil {
					fmt.Println(err)
					showNotification("Download failed: " + err.Error())
					return
				}
				downloadFunc(homeDir)
			}
		}, window)
	})

	// Create a horizontal box to hold widgets
	content := container.NewVBox(
		widget.NewLabel("Upload a file:"),
		container.NewHBox(
			widget.NewButton("Select File", func() {
				dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
					if err == nil && reader != nil {
						path := reader.URI().String()[5:] // removing "file://"
						filePathEntry.SetText(path)
					}
				}, window)
			}),

			widget.NewButton("Select Folder", func() {
				dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
					if err == nil && uri != nil {
						path := uri.Path()
						filePathEntry.SetText(path)
					}
				}, window)
			}),
		),
		filePathEntry,
		codeEntry,
		uploadButton,
		widget.NewLabel("Download a file:"),
		codeRecieve,
		downloadButton,
	)

	window.SetContent(content)
	window.ShowAndRun()
}
