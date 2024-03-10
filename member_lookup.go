package gohoa

import (
	"fmt"
	"strings"
)

type MemberLookup struct {
	dirQuerySvc *DirQueryService
	trieNum     *Trie
	trieStreet  *Trie
}

func NewMemberLookup() *MemberLookup {

	dirQueryService := NewDirQueryService()

	//Load all the static street mappings into memmory
	LoadStreetMappingsJson()

	//Load all members into memmory
	var members []Member
	dirQueryService.FindAllMembers(&members)

	trieNum := InitTrie(NUM_TRIE)
	trieText := InitTrie(TEXT_TRIE)

	//Add all members to the trie
	for _, member := range members {
		trieNum.Insert(fmt.Sprintf("%d", member.PAddress.Number))
	}

	streetMappings := GetStreetMap()
	for _, sm := range streetMappings {
		if sm.TrieNames != nil {
			for _, name := range sm.TrieNames {
				trieText.Insert(name)
			}
		} else {
			lowerName := strings.ToLower(sm.StreetName)
			trieText.Insert(lowerName)
		}
	}

	return &MemberLookup{dirQueryService, trieNum, trieText}
}

func (ml *MemberLookup) SuggestStreetName(prefix string) []StreetMapping {
	trieMatches := ml.trieStreet.Suggestions(prefix)
	var streetMappings []StreetMapping

	if trieMatches != nil {
		for _, match := range trieMatches {
			streetMappings = append(streetMappings, GetStreetMappingFromTrieName(match))
		}
	}
	return streetMappings

}

func (ml *MemberLookup) SuggestNumber(prefix string) []string {
	return ml.trieNum.Suggestions(prefix)
}

func (ml *MemberLookup) FindMembersByStreetNumber(number string) ([]PropertyAddress, error) {
	return ml.dirQuerySvc.FindMembersByStreetNumber(number)
}
