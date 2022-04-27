# standard filehandler

this is the standard file handler that is access via interface:

```
type FileHandler interface {
	ConvertMarkdown(filePath, pageTitle string, parentMetadata map[string]interface{}) (*FileContents, error)
	ConvertFolder(filePath, pageTitle string, parentMetadata map[string]interface{}) (*FileContents, error)
	ProcessOtherFile(filePath, pageTitle string, parentMetadata map[string]interface{}) (*FileContents, error) // for any other logic
}
```