package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/robododge/gohoa"

	"github.com/brianvoe/gofakeit/v6"
)

var (
	suffixes  = [...]string{"Street", "Court", "Drive", "Circle", "Trail", "Way"}
	suffixabb = [...]string{"st", "ct", "dr", "tr"}
)

const (
	JSON_SLIM_MEMBERS = "slim_members.json"
	JSON_SLIM_ANON    = "slim_members_anon.json"
	HOA_FILE          = "hoa_dir.json"
)

func main() {
	var anonymizeMode bool
	flag.BoolVar(&anonymizeMode, "anonymize", false, "Enable anonymizer mode")
	flag.Parse()
	if anonymizeMode {
		fmt.Println("Running in anonymizer mode")
		runAnonimization()
	} else {
		runCleaner()
	}
}

func runCleaner() {

	fmt.Println(" HOA directory cleaner")
	jsonFile, err := os.Open(HOA_FILE)
	if err != nil {
		log.Fatal("cannot open ", HOA_FILE)
	}
	defer jsonFile.Close()

	dec := json.NewDecoder(bufio.NewReaderSize(jsonFile, 1024))

	t, err := dec.Token()
	if err != nil {
		log.Fatal("decoder error")
	}
	fmt.Printf("%T: %v\n", t, t)

	var i int

	for dec.More() {
		i++

		if t, err := dec.Token(); err == nil {
			tv, ok := t.(string)
			if ok {
				if tv == "member" {
					if t, err := dec.Token(); err == nil {
						_, ok = t.(json.Delim)
						if ok {
							fmt.Println("Found the members array in interations:", i)
							stripNameAddress(dec)
							break

						}

					}
				}

			}

		}

	}

}

func stripNameAddress(dec *json.Decoder) {

	f, err := os.Create(JSON_SLIM_MEMBERS)
	if err != nil {
		log.Fatal("Cannot open file for writing members")
	}
	defer f.Close()
	err = writeSimple(f, "[")
	if err != nil {
		newErr := fmt.Errorf("WE got an eror writing initial part of file: %V ", err)
		fmt.Println(newErr)
	}

	member := &gohoa.Member{}
	mi := 0
	for {
		err := dec.Decode(member)
		if err != nil {
			fmt.Printf("Total members before exit: %d\n", mi)
			if tok, ert := dec.Token(); ert == nil {
				//if we error out because the members array is done
				if dl, ok := tok.(json.Delim); ok && dl.String() == "]" {
					fmt.Println("Successful end of members")
					_ = writeSimple(f, "]")
					break
				}
			}
			log.Fatal("Errror get'n member: ", err)
		}

		populateStreetDetails(member)

		if mi > 0 {
			_ = writeSimple(f, ",\n")
		}
		fmt.Println(*member)
		writeEntry(f, *member)
		mi++
	}

}

func populateStreetDetails(member *gohoa.Member) {

	cleanAddr := strings.TrimSpace(member.PAddress.Addr1)
	splits := strings.Fields(cleanAddr)

	if len(splits) >= 2 {

		member.PAddress.Number = parseStreetNum(splits[0], member)

		if len(splits) > 2 && notSuffix(splits[2]) {
			member.PAddress.StreetName = fmt.Sprintf("%s %s", splits[1], splits[2])
		} else {
			member.PAddress.StreetName = splits[1]
		}
	}

}

func parseStreetNum(snum string, member *gohoa.Member) int {
	if sval, err := strconv.ParseInt(snum, 10, 16); err == nil {
		return int(sval)
	} else {
		fmt.Printf("Oops! error for member %s, \n %s\n", member.Contacts[0].FirstName, err)
	}
	return 0
}

func notSuffix(val string) bool {

	sliceSuffix := suffixes[:]
	found := slices.Contains(sliceSuffix, val)

	if !found {
		lval := strings.ToLower(val)
		for _, abb := range suffixabb {
			found = strings.HasPrefix(lval, abb)
			if found {
				break
			}
		}
	}
	return !found
}

func writeEntry(file *os.File, m gohoa.Member) {
	if b, err := json.Marshal(m); err == nil {
		file.Write(b)
	}
}
func writeSimple(file *os.File, s string) error {
	_, err := file.Write([]byte(s))
	return err
}

func runAnonimization() {
	gofakeit.Seed(0)
	allMembers := gohoa.NewAllMembers()
	allMembers.PopulateFromJsonFile(JSON_SLIM_MEMBERS)

	hundredMemebers := allMembers.Members[:100]

	streetName := fmt.Sprintf("%s %s", gofakeit.StreetName(), gofakeit.StreetSuffix())
	streetNum := 1
	for i := range hundredMemebers {
		if streetNum%10 == 0 { //every 10th member, change the street
			streetName = fmt.Sprintf("%s %s", gofakeit.StreetName(), gofakeit.StreetSuffix())
			streetNum = streetNum + 1000
		}
		hundredMemebers[i].PAddress.Number = streetNum
		hundredMemebers[i].PAddress.StreetName = streetName
		hundredMemebers[i].PAddress.Addr1 = fmt.Sprintf("%d %s", streetNum, streetName)
		fakeLastName := gofakeit.LastName()
		for j := range hundredMemebers[i].Contacts {
			fakeFirstName := gofakeit.FirstName()
			hundredMemebers[i].Contacts[j].FirstName = fakeFirstName
			hundredMemebers[i].Contacts[j].LastName = fakeLastName
		}

		streetNum++
	}

	//sart the writing proceess
	f, err := os.Create(JSON_SLIM_ANON)
	if err != nil {
		log.Fatal("Cannot open file for writing anonymous members")
	}
	defer f.Close()
	err = writeSimple(f, "[")
	if err != nil {
		log.Fatal("Cannot write to anonymous members file")
	}

	for mi, member := range hundredMemebers {
		if mi > 0 {
			_ = writeSimple(f, ",\n")
		}
		fmt.Println(member)
		writeEntry(f, member)
	}
	_ = writeSimple(f, "]")
}
