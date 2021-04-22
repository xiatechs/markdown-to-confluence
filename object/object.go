// Package object is where to store the confluence variables
package object

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// ConfluenceVars is a struct to store the username/APIKey/Space variables to work with confluence API
type ConfluenceVars struct {
	ConfluenceUsernameEnv string
	ConfluenceAPIKeyEnv   string
	ConfluenceSpaceEnv    string
}

// ConfluenceObject that is used in the application
var ConfluenceObject = ConfluenceVars{}

// Save a ConfluenceVars obj down to json named 'confobject.json' located in same folder as the app
func (c ConfluenceVars) save() {
	item := &c

	output, err := json.MarshalIndent(item, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("confobject.json", output, 0600)
	if err != nil {
		fmt.Println(err)
	}
}

// Load a json into a ConfluenceVars obj
func (c *ConfluenceVars) Load() bool {
	item := *c
	jsonFile, _ := ioutil.ReadFile("confobject.json")

	err := json.Unmarshal(jsonFile, &item)
	if err != nil {
		fmt.Println(err)

		return false
	}

	if ok := item.isNotEmpty(); !ok {
		item.save()

		return false
	}

	*c = item

	return true
}

// isNotEmpty func determines whether the fields in the confobject.JSON are empty
func (c ConfluenceVars) isNotEmpty() bool {
	if c.ConfluenceUsernameEnv == "" || c.ConfluenceAPIKeyEnv == "" || c.ConfluenceSpaceEnv == "" {
		log.Println("Environment variable(s) not been set...")
		log.Printf("Username: %s", c.ConfluenceUsernameEnv)
		log.Printf("API KEY: %s", c.ConfluenceAPIKeyEnv)
		log.Printf("SPACE: %s", c.ConfluenceSpaceEnv)
		log.Println("Please update the confobject.json located wherever this application is located")

		return false
	}

	log.Printf("Username: %s", c.ConfluenceUsernameEnv)
	log.Printf("API KEY: %s", c.ConfluenceAPIKeyEnv)
	log.Printf("SPACE: %s", c.ConfluenceSpaceEnv)

	return true
}
