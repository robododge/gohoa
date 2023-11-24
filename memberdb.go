package gohoa

import (
	"fmt"
	"sort"
)

//Emulate a key/value store with a key lookup based on address
//This is a simple in-memory map, but could be extended to a database
//if the HOA has members in the thousands or more it would be wise to

func NewMemberDB() MemberDB {
	outDb := new(memberMapDB)
	outDb.mm = make(map[Mkey]*Member)
	//byStreet currently just an aggreatation of all members on a street
	outDb.byStreet = make(map[string][]*Member)
	return outDb
}

func (mdb *memberMapDB) AddConvience(member *Member) {
	key := MakeKey(member)
	mdb.Add(key, member)
}

func (mdb *memberMapDB) Add(key Mkey, member *Member) {
	mdb.byStreet[key.streetName] = append(mdb.byStreet[key.streetName], member)
	mdb.mm[key] = member
}

func MakeKey(member *Member) Mkey {
	mkey := Mkey{member.PAddress.StreetName, member.PAddress.Number}
	return mkey
}

func MakeKeyV(street string, number int) Mkey {
	mkey := Mkey{street, number}
	return mkey
}

func (mk Mkey) String() string {
	return fmt.Sprintf("Number:'%d' Street:'%s'", mk.streetNumber, mk.streetName)
}

func (mdb *memberMapDB) Fetch(key Mkey) (*Member, bool) {
	member, found := mdb.mm[key]
	return member, found
}

func (mdb *memberMapDB) Size() int {
	return len(mdb.mm)
}

// getAllMKeys is a backdoor getting all the keys in the map
func (mdb *memberMapDB) getAllMKeys() []Mkey {
	keys := make([]Mkey, 0, len(mdb.mm))
	for key := range mdb.mm {
		keys = append(keys, key)
	}
	return keys
}

// getAllMembersOnStreet is a backdoor getting all the members on a street
func (mdb *memberMapDB) getAllMembersOnStreet(street string) []*Member {

	if streetMembers, found := mdb.byStreet[street]; found {
		sort.Slice(streetMembers, func(i, j int) bool { return streetMembers[i].PAddress.Number < streetMembers[j].PAddress.Number })
		return streetMembers
	}

	return []*Member{}
}
