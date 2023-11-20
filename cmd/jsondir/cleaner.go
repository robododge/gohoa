package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

var (
	suffixes  = []string{"Street", "Court", "Drive", "Circle", "Trail", "Way"}
	suffixabb = []string{"st", "ct", "dr", "tr"}
)

func main() {

	fmt.Println(" HOA directory cleaner")
	jsonFile, err := os.Open("hoa_dir.json")
	if err != nil {
		log.Fatal("cannot open hoa_dir.json")
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

type Contact struct {
	RoleName  string `json:"role_name"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
}

type PropertyAddress struct {
	Addr1      string `json:"addr1"`
	StreetName string `json:",omitempty"`
	Number     int    `json:",omitempty"`
}

type Member struct {
	PAddress PropertyAddress `json:"prop_address"`
	Contacts []Contact       `json:"contact"`
}

type AllMembers struct {
	Members []Member
}

func stripNameAddress(dec *json.Decoder) {

	f, err := os.Create("slim_members.json")
	if err != nil {
		log.Fatal("Cannot open file for writing members")
	}
	defer f.Close()
	err = writeSimple(f, "[")
	if err != nil {
		newErr := fmt.Errorf("WE got an eror writing initial part of file: %V ", err)
		fmt.Println(newErr)
	}

	member := &Member{}
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

func populateStreetDetails(member *Member) {

	cleanAddr := strings.TrimSpace(member.PAddress.Addr1)
	splits := strings.Fields(cleanAddr)
	// for i, sp := range splits {
	// 	if i == 0 {
	// 	if sval, err := strconv.ParseInt(sp, 10, 16); err == nil {
	// 		member.PAddress.Number = int(sval)
	// 	} else {
	// 		member.PAddress.StreetName = fmt.Sprintf("%s %s")

	// 	}

	// }

	//Makesure to handle the case where mulitple spaces exist
	// var lSplits []string
	// for _, sraw := range splits {
	// 	if sraw != " " {
	// 		lSplits = append(lSplits, sraw)
	// 	}
	// }

	if len(splits) >= 2 {

		member.PAddress.Number = parseStreetNum(splits[0], member)

		// s1, s2 := strings.TrimSpace(splits[1]), strings.TrimSpace(splits[2])
		// s1, s2 := splits[1], splits[2]
		if len(splits) > 2 && notSuffix(splits[2]) {
			member.PAddress.StreetName = fmt.Sprintf("%s %s", splits[1], splits[2])
		} else {
			member.PAddress.StreetName = splits[1]
		}
	}

}

func parseStreetNum(snum string, member *Member) int {
	if sval, err := strconv.ParseInt(snum, 10, 16); err == nil {
		return int(sval)
	} else {
		fmt.Printf("Oops! error for member %s, \n %s\n", member.Contacts[0].FirstName, err)
	}
	return 0
}

func notSuffix(val string) bool {

	found := slices.Contains(suffixes, val)

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

func writeEntry(file *os.File, m Member) {
	if b, err := json.Marshal(m); err == nil {
		file.Write(b)
	}
}
func writeSimple(file *os.File, s string) error {
	_, err := file.Write([]byte(s))
	return err

}
