package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/robododge/gohoa"
)

const (
	MEMBERS_JSON = "../../sampledata/slim_members.json"
	ORDERS_CSV   = "../../sampledata/SpecialOrders.csv"
)

func main() {
	allMembers := gohoa.NewAllMembers()
	allMembers.PopulateFromJsonFile(MEMBERS_JSON)

	allOrders := gohoa.NewAllOrders()
	allOrders.LoadOrders(ORDERS_CSV)

	if err := allOrders.ValidateNumbers(); err != nil {
		log.Fatalf("Cannot proceed, found mis-match in entered street numbers %s\n", err.Error())
	}

	analysis := gohoa.NewAnalysis()
	analysis.LoadAllMembers(allMembers)
	matches, misses, err := analysis.CrossCheckOrders(allOrders)

	//sort the matches
	sort.Sort(gohoa.ByStreetNumber(matches))

	if err != nil {
		fmt.Println("Error found!")
	} else {
		fmt.Println("All addresses are clean!")
	}

	if len(matches) > 0 {
		t := table.NewWriter()
		t.SetTitle("Matches")

		for _, m := range matches {
			nameFound := fmt.Sprintf("%s %s", m.MatchedContact.FirstName, m.MatchedContact.LastName)
			t.AppendRow(table.Row{m.StreetNumber, m.StreetName, m.NeighborName, nameFound})
		}
		t.AppendHeader(table.Row{"Number", "Street", "Requestor Name", "Name In Directory"})
		fmt.Println(t.Render())
	}

	if len(misses) > 0 {
		t := table.NewWriter()
		t.SetTitle("Misses")
		for _, m := range misses {
			inhabitants := analysis.FetchInhabitants(m.StreetNumber, m.StreetName)
			var inhabitantsArry []string
			var inhabitantNames string
			for _, ih := range inhabitants {
				inhabitantsArry = append(inhabitantsArry, fmt.Sprintf("%s %s", ih.FirstName, ih.LastName))
			}
			if len(inhabitantsArry) > 0 {
				inhabitantNames = strings.Join(inhabitantsArry, ", ")
			}
			t.AppendRow(table.Row{m.StreetNumber, m.StreetName, m.Name, inhabitantNames})
		}
		t.AppendHeader(table.Row{"Number", "Street", "Requestor Name", "Actual Inhabitants"})
		fmt.Println(t.Render())
	}

	//Fetch all unique names
	uniqueStreets := analysis.UniqeStreetNames()
	sort.Strings(uniqueStreets)
	fmt.Println("Unique Streets:")
	for _, street := range uniqueStreets {
		fmt.Println(street)
	}

}
