package confluence

import (
	"log"

	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

var Iterator = 0

type Local struct {
	client *confluence.APIClient
}

func NewAPIClient() {
	client, err := confluence.CreateAPIClient()
	if err != nil {
		log.Println(err)
		return
	}

	l := &Local{
		client: client,
	}

	return
}

func (l *Local) CRUD(file *filehandler.FileContents, parentMetaData map[string]interface{}) (map[string]interface{}, error) {
	if file == nil {
		return nil, nil
	}

	var returnData map[string]interface{}

	switch file.MetaData["type"] {
	case "attachment":
		l.client.UploadAttachment(file.MetaData["type"])

		log.Println("CRUD", file.MetaData["type"], file.MetaData["title"], " PARENT:", "id:", parentMetaData["id"], "title:", parentMetaData["title"])

	case "markdown":
		Iterator++

		returnData = map[string]interface{}{
			"root":  parentMetaData["id"],
			"id":    Iterator,
			"title": file.MetaData["title"],
		}

		log.Println("CRUD", file.MetaData["type"], file.MetaData["title"], " PARENT:", "id:", parentMetaData["id"], "title:", parentMetaData["title"])
	case "folderpage":
		Iterator++

		returnData = map[string]interface{}{
			"root":  parentMetaData["id"],
			"id":    Iterator,
			"title": file.MetaData["title"],
		}

		log.Println("CRUD", file.MetaData["type"], file.MetaData["title"], " PARENT:", "id:", parentMetaData["id"], "title:", parentMetaData["title"])
	}

	return returnData, nil
}
