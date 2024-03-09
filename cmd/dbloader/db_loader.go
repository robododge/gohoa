package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/robododge/gohoa"
)

func main() {
	populatedNew := flag.Bool("new", false, "a bool")
	revalidate := flag.Bool("revalidate", false, "a bool")
	findcall := flag.Bool("findcall", false, "a bool")

	flag.Parse()

	flagsFound := []bool{*populatedNew, *revalidate, *findcall}
	flagsCount := 0
	for _, f := range flagsFound {
		if f {
			flagsCount++
		}
	}
	if flagsCount != 1 {
		log.Fatal("Can only have one flag")
	}
	if flagsCount == 0 {
		log.Fatal("Must have one flag")
	}

	gohoa.LoadStreetMappingsJson()

	if *populatedNew {
		gohoa.NewDirLoader().PopulateMongoFromJSON()
	} else if *revalidate {
		gohoa.NewDirLoader().RevalidateMongoFromJSON()
	} else if *findcall {
		dirSvc := gohoa.NewDirQueryService()
		dirSvc.FindCountByStreetName()
		var members []gohoa.Member
		dirSvc.FindAllMembers(&members)
		fmt.Println("Found ", len(members), " members")
	}

	// gohoa.NewDirLoader().PopulateMongoFromJson()
	// gohoa.NewDirLoader().RevalidateMongoFromJson()

	// gohoa.NewDirQueryService().FindCountByStreetName()

}
