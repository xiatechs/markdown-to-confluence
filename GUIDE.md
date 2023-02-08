v2.1
# markdown-to-confluence end user guide:
# if you want to adjust these rules, the tool needs to be adjusted.

```
- only the /docs folder of the repository will be copied across by default - set the flag ONLY-DOCS to "false" to copy all valid files
	- warning: if copying across all folders the hierarchy may not be replicated perfectly

- any folder where there are valid markdown files will be considered 'valid/alive' & contents mirrored to confluence

- the tool can generate the pages into an already existing confluence page (set PARENT-PAGE-ID to the pages ID) or it can be generated to a new root (set PARENT-PAGE-ID to 0)

- the tool will convert any markdown headers into Proper Case (so all words will start with an upper case and then be all lowercase thereafter)
- any circle brackets () in a markdown heading will be removed from the confluence page heading
	- so it is best to create the headings in markdown with Proper Case headings without any brackets in them
```
