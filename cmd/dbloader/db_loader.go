package main

import "github.com/robododge/gohoa"

func main() {

	// gohoa.NewDirLoader().PopulateMongoFromJson()
	// gohoa.NewDirLoader().RevalidateMongoFromJson()

	gohoa.NewDirQueryService().FindCountByStreetName()

}
