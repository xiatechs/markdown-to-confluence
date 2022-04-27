package standard

import (
	fh "github.com/xiatechs/markdown-to-confluence/filehandler"
)

// Basic - a basic filehandler that converts markdown files into HTML files, folders into 'generic folder HTML files'
// and images are passed back to API as metadata so the API can choose to handle them
type Basic struct {
}

func (l *Basic) ConvertMarkdown(filePath, pageTitle string, parentMetadata map[string]interface{}) (*fh.FileContents, error) {
	bytes, err := returnBytes(filePath)
	if err != nil {
		return nil, err
	}

	fileContents, err := convertMarkdownToHTML(bytes, parentMetadata)
	if err != nil {
		return nil, err
	}

	fileContents.MetaData["type"] = "markdown"
	fileContents.MetaData["filepath"] = filePath

	return fileContents, nil
}

// ConvertFolder - convert a Folder

func (l *Basic) ConvertFolder(filePath, pageTitle string, parentMetadata map[string]interface{}) (*fh.FileContents, error) {
	return &fh.FileContents{
		MetaData: map[string]interface{}{
			"type":     "folderpage",
			"title":    pageTitle,
			"filepath": filePath,
		},
		Body: []byte(filePath),
	}, nil
}

// ProcessOtherFile - process & convert other file types

func (l *Basic) ProcessOtherFile(filePath, pageTitle string, parentMetadata map[string]interface{}) (*fh.FileContents, error) {
	f := fh.NewFileContents()

	if isAcceptedImageFile(filePath) {
		f.MetaData["type"] = "attachment"
		f.MetaData["filepath"] = filePath
		f.MetaData["title"] = pageTitle
	}

	return f, nil
}
