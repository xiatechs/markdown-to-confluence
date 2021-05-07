# markdown-to-confluence/todo readme

## the todo package provides the logic necessary to gather TODO's in a repository codebase and store them

### The package contains one exported variable:
```
// this variable is of a type 'markdown.FileContents'.
MainTODOPage
```
### The package contains two exported functions:
```
// this function takes in a .go file and file name and collects all TODO's in the code
ParseGo(content []byte, filename string) 

// this function takes in the root directory of the repository and creates a page that can be posted on confluence.
GenerateTODO(rootDir string) *markdown.FileContents
```
