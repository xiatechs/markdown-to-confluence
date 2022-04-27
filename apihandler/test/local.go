package test

import (
	"log"

	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

var Iterator = 0

type Local struct{}

func (l *Local) CRUD(file *filehandler.FileContents, parentMetaData map[string]interface{}) (map[string]interface{}, error) {
	if file == nil {
		return nil, nil
	}

	var returnData map[string]interface{}

	switch file.MetaData["type"] {
	case "attachment":
		Iterator++

		returnData = map[string]interface{}{
			"root":  parentMetaData["id"],
			"id":    Iterator,
			"title": file.MetaData["title"],
		}

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
