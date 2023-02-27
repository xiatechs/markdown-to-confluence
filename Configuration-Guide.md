# Config guide / how to set up MTC.

If you don't know github actions and what they are - learn about them first here: https://docs.github.com/en/actions

1) Create an action YML - call it 'markdown.yml' or something similar and place it in the repo/.github/workflows/ folder in your repository:
```
on:
  push:
    branches: [master]
name: Markdown To Confluence
jobs:
  markdown:
    name: Markdown To Confluence Action
    runs-on: ubuntu-latest
    env: 
      SPACE: "XKB"
      PARENT-ROOT-ID: "123456789"
      PAGE-NAME: "{repo name}-GitHub-Docs"
      ONLY-DOCS: "true"
    steps:
      - name: gather branch details
        shell: bash
        run: |
             if [ -z "${GITHUB_HEAD_REF}" ]
             then
              echo NOT pull request, branch = $(echo ${GITHUB_REF#refs/heads/})
              echo "BRANCH_NAME=$(echo ${GITHUB_REF#refs/heads/})" >> $GITHUB_ENV
             else
              echo pull request, branch = $(echo ${GITHUB_HEAD_REF})
              echo "BRANCH_NAME=$(echo ${GITHUB_HEAD_REF})" >> $GITHUB_ENV
             fi
        id: extract_branch          
      - name: checkout markdown action
        uses: actions/checkout@v2-beta
        with:
          repository: xiatechs/markdown-to-confluence
          ref: refs/tags/v1.10
      - name: checkout branch
        uses: actions/checkout@v2
        with:
          ref: ${{ env.BRANCH_NAME }}
          path: "./${{ env.PAGE-NAME }}"
          fetch-depth: 0
      - name: run markdown to confluence action
        uses: ./
        with:
          key: ${{ secrets.CONFLUENCE_KEY }}
          space: "${{ env.SPACE }}"
          username: ${{ secrets.CONFLUENCE_USERNAME }}
          repo: "${{ env.PAGE-NAME }}"
          parentID: "${{ env.PARENT-ROOT-ID }}"
          url: "https://xiatech.atlassian.net"
          onlyDocs: "${{ env.ONLY-DOCS }}"

```

2) Edit the YML:
```
The bits you need to edit:

    env: 
      SPACE: "XKB"                #space is the name of the space in confluence you want the page to be in
      PARENT-ROOT-ID: "123456789" #parent root ID is the page ID of the root page to create the new mtc generated pages in (if 0 then pages don't get a parent and are generated to a root)
      PAGE-NAME: "{repo name} GitHub Docs"     #page name is the name of the page you want created in confluence (normally "{repo name}-Github-Docs") - CANNOT HAVE SPACES
      ONLY-DOCS: "true"              #only docs is a flag to decide whether it is only the /docs folder which will be copied to confluence (default should be true)
```

You can add tests/lint to the configuration if you want. 
