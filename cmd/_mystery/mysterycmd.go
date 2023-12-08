package main

// Trying to emulate a mystery where the last pointer
//  to the data in the underlying Map takes on
//  the pointer value for every entrie's value in the map
//https://go.dev/play/p/fVvLHbw99UD

import "fmt"

type Subscriber struct {
	Name    string
	Number  int
	Details string
}

type MyDB interface {
	Fetch(MyKey) (*Subscriber, bool)
	Add(MyKey, *Subscriber)
	AddConvience(*Subscriber)
}

type MyKey struct {
	Name   string
	Number int
}

type subMapDB struct {
	smap map[MyKey]*Subscriber
}

// Create an in-memory key database
// using pointers for values
func NewMyDB() MyDB {
	outDb := new(subMapDB)
	outDb.smap = make(map[MyKey]*Subscriber)
	return outDb
}

func (sdb *subMapDB) AddConvience(sub *Subscriber) {
	fmt.Printf("adding pointer1: %p \n", sub)
	key := MakeKey(sub)
	sdb.Add(key, sub)
}

func (sdb *subMapDB) Add(key MyKey, sub *Subscriber) {
	fmt.Printf("adding pointer2: %p \n", sub)
	sdb.smap[key] = sub
}

func MakeKey(sub *Subscriber) MyKey {
	mkey := MyKey{Name: sub.Name, Number: sub.Number}
	return mkey
}

func (mk MyKey) String() string {
	return fmt.Sprintf("Number:'%d' Name:'%s'", mk.Number, mk.Name)
}

func (sdb *subMapDB) Fetch(key MyKey) (*Subscriber, bool) {
	member, found := sdb.smap[key]
	return member, found
}

func (sdb *subMapDB) PrintPtrs() {
	for key, val := range sdb.smap {
		fmt.Printf("key:%s,  ptr:%p\n", key, val)
	}
}

func main() {
	subs := []Subscriber{
		{"Jack1", 1111, "one"},
		{"Jack2", 2222, "two"},
		{"Jack3", 3333, "three"},
		{"Jack4", 4444, "four"},
	}

	//When loading the database using a range loop, the loop variable is actually scoped as an outer variable to the for loop
	// thus the pointer value is in that variable is overwritten each time but the memory address where the new values are stored
	// doen't change, thus the same pointer.
	myDB := NewMyDB()
	for _, s := range subs {
		ss := s
		fmt.Printf("Outer loop pointer %p generated pointer: %p \n", &s, &ss)
		myDB.AddConvience(&ss)
	}

	search := MakeKey(&subs[1])

	subFound, found := myDB.Fetch(search)
	if found {
		fmt.Println("Sub found should be 'two', found:")
		fmt.Printf("pointer: %p\n", subFound)
	}

	s := myDB.(*subMapDB)
	s.PrintPtrs()

}
