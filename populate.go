package gohoa

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jszwec/csvutil"
)

func NewAllOrders() *AllOrders {
	return &AllOrders{orders: make([]Order, 0)}
}

func (ao *AllOrders) LoadOrders(filename string) {

	csvBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Cannot load csv file %s, error %s\n", filename, err.Error())
	}

	if err = csvutil.Unmarshal(csvBytes, &ao.orders); err != nil {
		fmt.Println("Error unmarshalling: ", err)
	}
}

func (ao *AllOrders) PrintOrders() {

	for i, o := range ao.orders {
		// fmt.Printf("Order[%d]: Name: %s,  number: %d\n", i, o.name, o.streetNumber)
		fmt.Printf("Order[%d]: %+v\n", i, o)
	}
}

func (ao *AllOrders) ValidateNumbers() error {

	for i, o := range ao.orders {
		if o.StreetNumber != o.StreetNumberV {
			return fmt.Errorf("order[%d] for %s has different street numbers %d vs %d", i, o.Name, o.StreetNumber, o.StreetNumberV)
		}
	}
	return nil
}

//TODO - add a new function for AllMembers

func NewAllMembers() *AllMembers {
	allM := AllMembers{make([]Member, 0)}
	return &allM
}

func (am *AllMembers) PopulateFromJsonFile(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Cannot open file %s eror %s\n", filename, err.Error())
	}
	defer file.Close()

	var members []Member

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&members); err != nil {
		log.Fatalf("Cannot unmarshall file %s into stucts error %s\n", filename, err.Error())
	}
	am.Members = members

}
