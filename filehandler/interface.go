package filehandler

// FileContents is a generic object that contains byte data
// and meta data of the file (can be anything)
type FileContents struct {
	MetaData map[string]interface{}
	Body     []byte
}

// NewFileContents creates a new filecontents obj
func NewFileContents() *FileContents {
	return &FileContents{
		MetaData: make(map[string]interface{}),
	}
}

//go:generate mockgen -destination=./filehandler_mocks.go -package=filehandler -source=interface.go

// FileHandler - to decouple handling of files from the API 
type FileHandler interface {
	ConvertMarkdown(filePath, pageTitle string, 
		parentMetadata map[string]interface{}) (*FileContents, error)
	ConvertFolder(filePath, pageTitle string, 
		parentMetadata map[string]interface{}) (*FileContents, error)
	ProcessOtherFile(filePath, pageTitle string, 
		parentMetadata map[string]interface{}) (*FileContents, error) // for any other logic
}
