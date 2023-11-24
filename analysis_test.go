package gohoa

import (
	"reflect"
	"testing"
)

var (
	members = []Member{
		{PAddress: PropertyAddress{
			Number: 5566, StreetName: "Westlane"},
			Contacts: []Contact{{RoleName: "Owner", FirstName: "Jim", LastName: "Dill"}}},
		{PAddress: PropertyAddress{ //This house hold has tons of occupants, would never fly in any HOA
			Number: 4567, StreetName: "Hawkins"},
			Contacts: []Contact{
				{RoleName: "Owner", FirstName: "Jan", LastName: "Hill"},
				{RoleName: "Tenant", FirstName: "John", LastName: "Smith"},
				{RoleName: "Owner", FirstName: "Sarah", LastName: "Johnson"},
				{RoleName: "Tenant", FirstName: "Emily", LastName: "Davis"},
				{RoleName: "Owner", FirstName: "Michael", LastName: "Brown"},
				{RoleName: "Tenant", FirstName: "Jessica", LastName: "Wilson"},
				{RoleName: "Tenant", FirstName: "Emily", LastName: "AndersonAndersonAndersonAndersonAnderson"},
				{RoleName: "Tenant", FirstName: "Sophia", LastName: "BrownBrownBrownBrownBrownBrownBrownBrownBrownBrown"},
				{RoleName: "Tenant", FirstName: "Olivia", LastName: "DavisDavisDavisDavisDavisDavisDavisDavisDavisDavis"},
				{RoleName: "Tenant", FirstName: "Ava", LastName: "GarciaGarciaGarciaGarciaGarciaGarciaGarciaGarciaGarcia"},
				{RoleName: "Tenant", FirstName: "Mia", LastName: "HarrisHarrisHarrisHarrisHarrisHarrisHarrisHarrisHarrisHarris"},
				{RoleName: "Tenant", FirstName: "Charlotte", LastName: "JacksonJacksonJacksonJacksonJacksonJacksonJacksonJacksonJacksonJackson"},
				{RoleName: "Tenant", FirstName: "Amelia", LastName: "JohnsonJohnsonJohnsonJohnsonJohnsonJohnsonJohnsonJohnsonJohnsonJohnson"},
				{RoleName: "Tenant", FirstName: "Harper", LastName: "MartinezMartinezMartinezMartinezMartinezMartinezMartinezMartinezMartinezMartinez"},
				{RoleName: "Tenant", FirstName: "Evelyn", LastName: "MillerMillerMillerMillerMillerMillerMillerMillerMillerMiller"},
				{RoleName: "Tenant", FirstName: "Abigail", LastName: "RobinsonRobinsonRobinsonRobinsonRobinsonRobinsonRobinsonRobinsonRobinsonRobinson"},
			}},
	}
)

func Test_hasFuzzyContactMatch(t *testing.T) {
	type args struct {
		memberf *Member
		first   string
		last    string
	}
	tests := []struct {
		name         string
		args         args
		wantContactf Contact
		wantFound    bool
	}{
		{
			"Test exact match",
			args{&members[0], "jim", "dill"},
			members[0].Contacts[0],
			true,
		},
		{
			"Find some fuzzy match",
			args{&members[1], "Emily", "AndAndersonAndersonAnderson"},
			members[1].Contacts[6],
			true,
		},
		{
			"Best fuzzy match from group",
			args{&members[1], "jess", "Wi"},
			members[1].Contacts[5],
			true,
		},
		{
			"Best fuzzy match for charlotte",
			args{&members[1], "cha", "jacksonjacksonjackson"},
			members[1].Contacts[11],
			true,
		},
		{
			"A non-match, Charolette",
			args{&members[1], "char", "levinstien"},
			members[1].Contacts[11],
			false,
		},
		{
			"A non-match, Olivia",
			args{&members[1], "Olivia", "levinstien"},
			members[1].Contacts[8],
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContactf, gotFound := hasFuzzyContactMatch(tt.args.memberf, tt.args.first, tt.args.last)
			if gotFound != tt.wantFound {
				t.Errorf("hasFuzzyContactMatch() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
			if tt.wantFound && !reflect.DeepEqual(gotContactf, tt.wantContactf) {
				t.Errorf("hasFuzzyContactMatch() gotContactf = %v, want %v Tested( %s, %s)", gotContactf, tt.wantContactf, tt.args.first, tt.args.last)
			}
		})
	}
}
