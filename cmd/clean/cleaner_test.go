package main

import (
	"testing"

	"github.com/robododge/gohoa"
)

func Test_populateStreetDetails(t *testing.T) {
	type args struct {
		member       *gohoa.Member
		resultNumber int
		resultStName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Space before streed num",
			args: args{
				member:       &gohoa.Member{PAddress: gohoa.PropertyAddress{Addr1: " 4573 Thorny Court"}, Contacts: []gohoa.Contact{{RoleName: "Owner", FirstName: "David", LastName: "Louve"}}},
				resultNumber: 4573,
				resultStName: "Thorny",
			},
		},
		{
			name: "Double space in name",
			args: args{
				member:       &gohoa.Member{PAddress: gohoa.PropertyAddress{Addr1: " 4577  Lancey"}, Contacts: []gohoa.Contact{{RoleName: "Owner", FirstName: "David", LastName: "Louve"}}},
				resultNumber: 4577,
				resultStName: "Lancey",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			populateStreetDetails(tt.args.member)
			if tt.args.resultNumber != tt.args.member.PAddress.Number {
				t.Errorf("St Number expected %d instead %d", tt.args.resultNumber, tt.args.member.PAddress.Number)
			}
			if tt.args.resultStName != tt.args.member.PAddress.StreetName {
				t.Errorf("St Name expected '%s' instead '%s'", tt.args.resultStName, tt.args.member.PAddress.StreetName)
			}
		})
	}
}
