package gohoa

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	suffixes  = [...]string{"street", "court", "drive", "circle", "trail", "way"}
	suffixabb = [...]string{"st", "ct", "dr", "tr"}
)

/*
A simple implmentation of this address parser will work with "your" HOA address

	it is not feasible to provide an all-encompassing addresss pareser, the
	logistics of doing so are documented here: https://www.reddit.com/r/golang/comments/i2jo1o/comment/g05f9rq/?utm_source=share&utm_medium=web2x&context=3
*/
type AddressLineParser interface {
	ParseStreetSegments(streetName string, tokCount int) []string
}

var _ AddressLineParser = (*segmentBasedStreetParser)(nil)

type segmentBasedStreetParser struct{}

// ParseAddressSegments - parse up to the max segments of a simple street name
func (sp *segmentBasedStreetParser) ParseStreetSegments(streetName string, maxSegments int) []string {
	splits := strings.Fields(streetName)
	var segments []string

	loopCnt := len(splits)
	if loopCnt > maxSegments {
		loopCnt = maxSegments
	}

	for i := 0; i < loopCnt; i++ {
		if notSuffix(splits[i]) {
			seg := segmentTransformer(splits[i])
			segments = append(segments, seg)
		} else {
			break
		}
	}
	return segments
}

func notSuffix(val string) bool {

	sliceSuffix := suffixes[:]
	lval := strings.ToLower(val)
	found := slices.Contains(sliceSuffix, lval)

	if !found {
		for _, abb := range suffixabb {
			found = strings.HasPrefix(lval, abb)
			if found {
				break
			}
		}
	}
	return !found
}

func segmentTransformer(seg string) string {
	if strings.ToLower(seg) == "saint" {
		return "St."
	}
	return seg
}

// SORT for the order matches
func (a ByStreetNumber) Len() int           { return len(a) }
func (a ByStreetNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStreetNumber) Less(i, j int) bool { return a[i].StreetNumber < a[j].StreetNumber }

// Translate streets into 3 character identifiers

func getStreetID(street string) string {
	switch street {
	case "street":
		return "st."
	case "court":
		return "ct."
	case "drive":
		return "dr."
	case "circle":
		return "cir."
	case "trail":
		return "trl."
	case "way":
		return "wy."
	default:
		return ""
	}
}

type AddressParseError struct {
	StreetNumber int
	StreetName   string
	Err          error
}

func (e *AddressParseError) Error() string {
	return fmt.Sprintf("Could not use number: '%d' street: %s as address", e.StreetNumber, e.StreetName)
}

// String returns a human-readable error message.
// func (e *AddressParseError) String() string {
// 	return fmt.Sprintf("Could not use number: '%d' street: %s as address", e.StreetNumber, e.StreetName)
// }

func CreateMongoIDForDiretory(m *Member) (string, string, error) {
	sName := GetShortIDFromStreetName(m.PAddress.StreetName)
	if sName == "" || m.PAddress.Number == 0 {
		return "", sName, &AddressParseError{m.PAddress.Number, m.PAddress.StreetName, errors.New("address error")}
	}
	outID := fmt.Sprintf("%s-%d-%d", sName, m.PAddress.Number, m.MemberId)
	return outID, sName, nil
}
