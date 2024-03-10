package gohoa

type Contact struct {
	RoleName  string `json:"role_name"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
}

type PropertyAddress struct {
	AddressID       int    `json:"address_id,string"`
	Addr1           string `json:"addr1"`
	StreetName      string `json:"StreetName"`
	Number          int    `json:"Number"`
	StreetNameShort string `json:"StreetNameShort"`
	StreetNameClean string `json:"StreetNameClean"`
}

type Member struct {
	ID         string          `json:"_id,omitempty" bson:"_id"`
	MemberId   int32           `json:"member_id,string"`
	MemberName string          `json:"member_name"`
	OrderId    int32           `json:"order_id,string"`
	PAddress   PropertyAddress `json:"prop_address"`
	Contacts   []Contact       `json:"contact"`
}

type AllMembers struct {
	Members []Member
}

/* Order types */
type Order struct {
	Name          string `csv:"Name"`
	Email         string `csv:"Email"`
	StreetName    string `csv:"Street Name"`
	StreetNumber  int    `csv:"Street Number"`
	StreetNumberV int    `csv:"Street Number (verifiy)"`
}

type AllOrders struct {
	orders []Order
}

/*  Data mapping for searching addresses in-memory **/
type Mkey struct {
	streetName   string
	streetNumber int
}

type MemberDB interface {
	Fetch(key Mkey) (*Member, bool)
	Add(key Mkey, member *Member)
	AddConvience(member *Member)
	Size() int
	AllMembers() []Member
}

type memberMapDB struct {
	mm       map[Mkey]*Member
	byStreet map[string][]*Member
}

type OrderMatch struct {
	NeighborName     string
	StreetNumber     int
	StreetName       string
	DiretoryContacts []Contact
	MatchedContact   Contact
}

type OrderMiss Order

type ByStreetNumber []OrderMatch
