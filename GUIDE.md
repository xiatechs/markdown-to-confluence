v1.1
# markdown-to-confluence end user guide:
# if you want to adjust these rules, the tool needs to be adjusted.

```
1) URL links i.e [link](https://hey.com) can be be one instance per row. 
	- the tool currently cannot parse multiple URL links one one line

2) markdown files should be within same folder as the code that they are talking about
	- this enables the tool to also generate plantuml diagrams for the code.

3) any images you want to link using relative links in the markdown must be saved in the same folder as the markdown files
	- alternatively, the image must be linked using an absolute path. Then it can be saved anywhere.

4) if you want to use those images in the markdown file, it must be linked locally i.e ![Diagram](diagram.png):
	- alternatively, as previously mentioned, you can just use absolute paths i.e:
	![Diagram](https://absolute-url/diagram.png)
  
5) you can have as many markdown files as you want in one folder, but typically
it should just be one for each package.
	- the markdown tool will create a page for each folder, obviously because
	if multiple markdown pages are in the same folder you need to be able to
	place them in the same page.
```
