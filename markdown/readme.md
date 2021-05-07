# markdown-to-confluence/markdown readme

## the markdown package is to enable working with and parsing markdown documents

### The package contains one exported struct:
```
// FileContents contains information from a file after being parsed from markdown.
FileContents{}
```

### The package contains one exported function:
```
// ParseMarkdown is a function that uses external parsing library to grab markdown contents
// and return a filecontents object (a page to be uploaded to confluence wiki)
ParseMarkdown(rootID int, content []byte) (*FileContents, error)
```
