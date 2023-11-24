package gohoa

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/avito-tech/normalize"
)

const (
	DEFAULT_ADDR_SEG_LEN = 2
	SIM_THRESHOLD        = 0.7
)

type Analysis struct {
	db     MemberDB
	ap     AddressLineParser
	segLen int
}

type rankContact struct {
	contact  Contact
	distance int
	distF    int
	distL    int
}

func NewAnalysis() *Analysis {
	a := Analysis{}
	a.db = NewMemberDB()
	a.ap = &segmentBasedStreetParser{}
	a.segLen = DEFAULT_ADDR_SEG_LEN

	lenSeg, found := os.LookupEnv("ADDR_SEG_LENGTH")
	if found {
		if parsedSeg, err := strconv.ParseInt(lenSeg, 10, 16); err == nil {
			a.segLen = int(parsedSeg)
		}
	}
	return &a
}

func (a *Analysis) LoadAllMembers(allMembers *AllMembers) {
	for _, m := range allMembers.Members {
		mTemp := m
		a.db.Add(MakeKey(&m), &mTemp)
	}
}

func (a *Analysis) CrossCheckOrders(ao *AllOrders) ([]OrderMatch, []OrderMiss, error) {
	var matches []OrderMatch
	var misses []OrderMiss
	for _, o := range ao.orders {
		streetClean := a.parseStreetName(o.StreetName)
		mKey := Mkey{streetClean, o.StreetNumber}
		// fmt.Printf("order name: %s\n", o.Name)
		m, found := a.db.Fetch(mKey)

		// fmt.Printf("Fetched: %s %s\n", m.Contacts[0].LastName, m.Contacts[0].FirstName)
		if !found {
			fmt.Printf("Could not find Mkey: %v\n", mKey)
			misses = append(misses, OrderMiss(o))
		} else {
			// fmt.Printf("Appendig match: %s %s\n", m.Contacts[0].LastName, m.Contacts[0].FirstName)
			//Must test that the contact actually matches as well
			contactFound, nameMatches := MatchFuzzyNameFromOrder(m, o)
			if nameMatches {
				om := OrderMatch{NeighborName: o.Name, StreetNumber: o.StreetNumber, StreetName: o.StreetName,
					DiretoryContacts: m.Contacts, MatchedContact: contactFound}
				matches = append(matches, om)
			} else {
				fmt.Printf("Could not find name from order in contact list: %v\n", o.Name)
				misses = append(misses, OrderMiss(o))
			}
		}

	}

	var err error
	if len(misses) > 0 {
		err = fmt.Errorf("found %d misses", len(misses))
	}

	return matches, misses, err
}

func (a *Analysis) FetchInhabitants(streetNum int, streetName string) []Contact {
	streetClean := a.parseStreetName(streetName)
	mKey := Mkey{streetClean, streetNum}
	m, found := a.db.Fetch(mKey)
	if !found {
		fmt.Printf("Could not find Mkey: %v\n", mKey)
		return nil
	}
	return m.Contacts
}

func (a *Analysis) parseStreetName(rawStreetName string) string {
	streetName := strings.Join(a.ap.ParseStreetSegments(rawStreetName, a.segLen), " ")
	return streetName
}

func MatchFuzzyNameFromOrder(m *Member, o Order) (Contact, bool) {
	//Parse the name coming as a single string
	tokens := strings.Fields(o.Name)

	if len(tokens) == 2 {
		return hasFuzzyContactMatch(m, tokens[0], tokens[1])
	}

	return Contact{}, false

}

// check if the member given has any contact matching the name
func hasFuzzyContactMatch(memberf *Member, first string, last string) (contactf Contact, found bool) {
	//Initial problem with this fuzzy library is each of the characters must be present as a subset of
	// characters in the target string.   If the compared string has a character that doesn't exist
	// in the target, then that makes the entire string be a no-match.  Kinda like a bloom filter I guess

	//first check the more strict version of Match
	for _, contact := range memberf.Contacts {
		// if fuzzy.MatchFold(last, contact.LastName) && fuzzy.MatchFold(first, contact.FirstName) {
		if last == contact.LastName && first == contact.FirstName {
			contactf, found = contact, true
			return
		}
	}

	var rc []rankContact

	//next try Levneshtein distance backed algorithm
	for _, contact := range memberf.Contacts {

		// concatTarget := fmt.Sprintf("%s%s", contact.FirstName, contact.LastName)
		// concatInput := fmt.Sprintf("%s%s", first, last)

		if similarFirst, similarLast := normalize.AreStringsSimilar(contact.FirstName, first, SIM_THRESHOLD),
			normalize.AreStringsSimilar(contact.LastName, last, SIM_THRESHOLD); similarLast && similarFirst {
			targetNorm := normalize.Many([]string{contact.FirstName, contact.LastName})
			inputNorm := normalize.Many([]string{first, last})
			distF := levenshtein.ComputeDistance(targetNorm[0], inputNorm[0])
			distL := levenshtein.ComputeDistance(targetNorm[1], inputNorm[1])
			rc = append(rc, rankContact{contact, distF + distL, distF, distL})
		}

	}

	//Sort if needed, return the best
	if len(rc) > 0 {
		if len(rc) > 1 {
			sort.Slice(rc, func(i, j int) bool {
				return rc[i].distance < rc[j].distance
			})
		}
		contactf, found = rc[0].contact, true
	}

	return
}

func (a *Analysis) FetchSingle(streetNum int, streetName string) (*Member, bool) {

	return a.db.Fetch(MakeKeyV(a.parseStreetName(streetName), streetNum))

}

func (a *Analysis) UniqeStreetNames() []string {

	var uniqueNames []string
	keys := a.db.(*memberMapDB).getAllMKeys()
	unique := make(map[string]bool)
	for _, key := range keys {
		unique[key.streetName] = true
	}
	for k := range unique {
		uniqueNames = append(uniqueNames, k)
	}

	return uniqueNames

}

func (a *Analysis) FetchAllMembersOnStreet(streetName string) []*Member {
	return a.db.(*memberMapDB).getAllMembersOnStreet(a.parseStreetName(streetName))
}
