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
	fileContents.MetaData["title"] = pageTitle

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
		Body: []byte(`<p>Welcome to the '<b>` + pageTitle + `</b>' folder of this Xiatech code repo.</p>
		<p>This folder full path in the repo is: ` + filePath + `</p>
<p>You will find attachments/images for this folder via the ellipsis at the top right.</p>
<p>Any markdown or subfolders is available in children pages under this page.</p>`),
	}, nil
}

// ProcessOtherFile - process & convert other file types

func (l *Basic) ProcessOtherFile(filePath, pageTitle string, parentMetadata map[string]interface{}) (*fh.FileContents, error) {
	f := fh.NewFileContents()

	if isAcceptedImageFile(filePath) {
		f.MetaData["type"] = "image"
		f.MetaData["filepath"] = filePath
		f.MetaData["title"] = pageTitle
	}

	if isGoFile(filePath) {
		f.MetaData["type"] = "go"
		f.MetaData["filepath"] = filePath
		f.MetaData["title"] = pageTitle
	}

	return f, nil
}
