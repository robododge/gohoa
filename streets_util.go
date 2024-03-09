package gohoa

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

// Gloabl variable, in memory at all times for street lookup
var streetNameToID = map[string]string{"Contrivy": "CTV", "Central": "CTL"}

// This mapping shouild be stored in the streets.json file
var shortIdToStreetMapping = map[string]StreetMapping{
	"CTV": {StreetName: "Contrivy", ShortID: "CTV", TrieNames: []string{"contrivy"}, FullName: "Contrivy Rd."},
	"CTL": {StreetName: "Central", ShortID: "CTL", TrieNames: []string{"central"}, FullName: "Central Ave."},
}

var trieNameMap = map[string]StreetMapping{}

var loadOnce sync.Once

type StreetMapping struct {
	StreetName string   `json:"name"`
	FullName   string   `json:"fullname"`
	ShortID    string   `json:"code"`
	TrieNames  []string `json:"trienames"`
}

func LoadStreetMappingsJson() {

	loadOnce.Do(func() {

		if GetConfig().StreetsJson == "" {
			log.Println("Didn't find streets.json file, not loading from disk")
			return
		}

		file, err := os.Open(GetConfig().StreetsJson)
		if err != nil {
			log.Fatalf("Cannot open streets json file %s eror %s\n", GetConfig().StreetsJson, err.Error())
		}
		defer file.Close()

		var streetMappings []StreetMapping

		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&streetMappings); err != nil {
			log.Fatalf("Cannot unmarshall file %s into stucts error %s\n", GetConfig().StreetsJson, err.Error())
		} else {
			for _, sm := range streetMappings {
				streetNameToID[sm.StreetName] = sm.ShortID
				shortIdToStreetMapping[sm.ShortID] = sm

				//Populate the trieNameMap so full street info retrieval from trie is possible
				if sm.TrieNames != nil {
					for _, name := range sm.TrieNames {
						trieNameMap[name] = sm
					}
				} else {
					lowerName := strings.ToLower(sm.StreetName)
					trieNameMap[lowerName] = sm
				}
			}
		}
	})
}

func GetShortIdFromStreetName(streetName string) string {
	return streetNameToID[streetName]
}
func GetStreetNameFromShortId(shortId string) StreetMapping {
	return shortIdToStreetMapping[shortId]
}

func GetStreetMap() map[string]StreetMapping {
	return shortIdToStreetMapping
}

func GetStreetMappingFromTrieName(trieName string) StreetMapping {
	return trieNameMap[trieName]
}

func GetAllStreetMappings() []StreetMapping {
	var streetMappings []StreetMapping
	for _, sm := range shortIdToStreetMapping {
		streetMappings = append(streetMappings, sm)
	}
	return streetMappings
}
