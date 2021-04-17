# notionbackup

A small utility command line application that can recursively download Notion pages.

I needed something scriptable that could periodically download my whole Notion workspace for backup purposes.

I have written a small blog post about this tool if you want to check it out: [Automate your Notion backups](https://blog.shayan.sx/automate-your-notion-backups.html)

## How To Use

Build it yourself or download a pre-built binary from the [Releases page](https://github.com/5hay/notionbackup/releases).

Set the following env variables

- `NOTION_TOKEN` (the token_v2 cookie value, just google "notion token_v2")
- `NOTION_PAGEID` (the id of the page you want to download recursively)
- _Optional_ `NOTION_EXPORTDIR` (the folder where the created .zip file should be placed in, **defaults to the current directory**)
  - Only specify the directory, the filename will be created for you
- _Optional_ `NOTION_EXPORTTYPE` ("html" or "markdown", **defaults to markdown**)

Now you can just run `./notionbackup`.

## Building

Clone this repository and then run `go build -o notionbackup -ldflags="-s -w" app.go`

## Special Thanks

- kjk for [notionapi](https://github.com/kjk/notionapi)
