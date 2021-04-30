# Markdown to Confluence Action

This Action will trawl through a repository on github & create a tree of pages for all markdown / images in a pattern like below:

```
Repo folder (attachments)
- markdown1
- markdown2

-- subfolder (attachments)
-- markdown1
-- markdown2

```

Any pages that are edited, and pushed, will be picked up by the action & automatically updated on confluence.

This uses the [Confluence REST API](https://developer.atlassian.com/cloud/confluence/rest/intro/)

requirements to run the script:
 - please provide the following env vars:
    "INPUT_CONFLUENCE_USERNAME"
    "INPUT_CONFLUENCE_API_KEY"
    "INPUT_CONFLUENCE_SPACE"
    "PROJECT_PATH"
   

