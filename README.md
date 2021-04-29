+++
title = "Markdown Action"
categories = ["Development", "Github Actions"]
date = "2021-03-10"
description = "A guide on how to use the markdown to confluence action"
slug = "markdown-to-confluence-guide"
+++


# Markdown to Confluence Action

This Action will find markdown files in a repository and read them, if they have an approriate [Front Matter](https://gohugo.io/content-management/front-matter/), it will create or update relevant pages in confluence.

This Action will also generate a plaintext uml & diagram of the codebase & upload them to the page.

This uses the [Confluence REST API](https://developer.atlassian.com/cloud/confluence/rest/intro/)

requirements to run the script:
 - please provide the following env vars:
    "INPUT_CONFLUENCE_USERNAME"
    "INPUT_CONFLUENCE_API_KEY"
    "INPUT_CONFLUENCE_SPACE"
    "PROJECT_PATH"
   

