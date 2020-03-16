package textloader

import (
	"io/ioutil"

	"github.com/wailsapp/wails"
)

// TextLoader is an utility for reading text files.
type TextLoader struct {
	Runtime *wails.Runtime
}

// NewTextLoader returns a new text loader.
func NewTextLoader() *TextLoader {
	return &TextLoader{}
}

// WailsInit initialize the structure for Wails
func (loader *TextLoader) WailsInit(r *wails.Runtime) error {
	loader.Runtime = r
	return nil
}

// GetTextFileContent let user select a text file and returns its content
// [output] string: the text content if no error
// otherwise the error code:
// #ErrorReadingFile: the file could not be read
func (loader *TextLoader) GetTextFileContent() string {
	selectedFile := loader.Runtime.Dialog.SelectFile()

	data, err := ioutil.ReadFile(selectedFile)
	if err != nil {
		return "#ErrorReadingFile"
	}

	result := string(data)

	return result
}

// GetTextFileContentManually let you read content from a text file without file chooser.
// path string: the absolute path
// [output] string: the text content if no error
// otherwise the error code:
// #ErrorReadingFile: the file could not be read
func (loader TextLoader) GetTextFileContentManually(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "#ErrorReadingFile"
	}

	result := string(data)

	return result
}
