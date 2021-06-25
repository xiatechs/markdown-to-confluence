# markdown-to-confluence/confluence readme

## the confluence package is for confluence api logic

### The package contains multiple exported structs:
```
// PageResults contains a slice of returned Pages
PageResults{}

//Page holds the data of a page from confluence wiki
Page{}

there are more exported structs within Page available in confluence/structs.go of this repo

// APIClient struct is for interacting with confluence
APIClient{}
```

### The package contains multiple exported methods:
```
APIClient{}

// CreateAPIClient creates the API client with relevant login details for confluence's API
CreateAPIClient() (*APIClient, error)
The username, API key and confluence space are provided via environment variables

// CreatePage is for creating a confluence page, root is the parent page ID
// if isroot is true the page is created as a root parent page in confluence
// it returns the page ID generated for the page
CreatePage(root int, contents *markdown.FileContents, isroot bool) (int, error)

// DeletePage deletes a confluence page using pageID to identify page to delete
DeletePage(pageID int) error

// UpdatePage updates a confluence page using pageID to identify page
UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents) error

// FindPage in confluence using title and returns page results
// if many is set to true it will also return the children pages of the page
FindPage(title string, many bool) (*PageResults, error)

// UploadAttachment to a page identified by page ID
UploadAttachment(filename string, id int) error
```
