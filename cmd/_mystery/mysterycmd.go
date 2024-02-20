package main

// Trying to emulate a mystery where the last pointer
//  to the data in the underlying Map takes on
//  the pointer value for every entrie's value in the map
//https://go.dev/play/p/fVvLHbw99UD

import (
	"fmt"
	"time"
)

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

	fmt.Print("Playing with select on channels\n")
	playWithSelectOnChannels()

}

func playWithSelectOnChannels() {

	var chan1 = make(chan int)
	var chan2 = make(chan int)
	var chan3 = make(chan string, 1)
	var chan4 = make(chan string, 1)

	go channelWatcher(chan1, chan2, chan3, chan4)

	go func() {
		time.Sleep(1 * time.Second)
		chan3 <- "hello"
	}()

	for i := 0; i < 10; i++ {
		chan2 <- i
		if i%2 == 0 {
			chan1 <- i
		}
	}
	out := <-chan4

	fmt.Printf("\nDone: %s\n", out)
}

func channelWatcher(chan1, chan2 chan int, chan3 chan string, chan4 chan<- string) {

	for {
		select {
		case <-time.Tick(1 * time.Second):
			fmt.Println("-- tick is done")
			chan4 <- "done-tick"
		case v1 := <-chan1:
			fmt.Printf("chan1: %d\n", v1)
		case v2 := <-chan2:
			fmt.Printf("chan2: %d\n", v2)
		case res := <-chan3:
			fmt.Printf("chan3 string result: %s\n", res)
			chan4 <- "done-chan3"
		}

	}

}
