# Markdown to Confluence Action

- This Action will crawl through a repository on github and create a tree of pages on confluence for all markdown / images:

![Diagram of action methodology](node.png)

## When this action is run:

1) A confluence page will be created containing the markdown documentation of a repository
2) Any new markdown pages in the repository will be uploaded or updated
3) A plaintext markup page will be generated and uploaded for code in pages with markdown
4) Old markdown pages online that have been removed in the repository will be deleted

## Features:

1) Folders with no content will be skipped to prevent a long chain of child pages & general confusion.
2) Images will be displayed in markdown pages - but ONLY if the images are stored in the same folder as the markdown page.

## Important:

1) There must be at least one markdown file in the root repository i.e README.md
2) All markdown must have a unique title. If two markdown have same the title - one will be ignored.
3) Markdown title is parsed by the first # header, or if that doesn't exist, the first ## or ### header.
4) Alternatively, title can be grabbed via TOML frontmatter.

## This action uses the [Confluence REST API](https://developer.atlassian.com/cloud/confluence/rest/intro/)

## Requirements to run the script:
 - please provide the following env vars:
    "INPUT_CONFLUENCE_USERNAME"
    "INPUT_CONFLUENCE_API_KEY"
    "INPUT_CONFLUENCE_SPACE"
    "PROJECT_PATH"
   

