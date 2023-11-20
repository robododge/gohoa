package main

import "testing"

func Test_populateStreetDetails(t *testing.T) {
	type args struct {
		member       *Member
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
				member:       &Member{PAddress: PropertyAddress{Addr1: " 4573 Thorny Court"}, Contacts: []Contact{{RoleName: "Owner", FirstName: "David", LastName: "Louve"}}},
				resultNumber: 4573,
				resultStName: "Thorny",
			},
		},
		{
			name: "Double space in name",
			args: args{
				member:       &Member{PAddress: PropertyAddress{Addr1: " 4577  Lancey"}, Contacts: []Contact{{RoleName: "Owner", FirstName: "David", LastName: "Louve"}}},
				resultNumber: 4577,
				resultStName: "Lancey",
			},
		},

		// {"prop_address":{"addr1":"4573  Lancelot","StreetName":" Lancelot","Number":4573},"contact":[{"role_name":"Owner","fname":"Jayme","lname":"Parks"},{"role_name":"Owner","fname":"Michael","lname":"Parks"}]},

		//{"prop_address":{"addr1":" 4573 Turnberry Court","StreetName":"4573 Turnberry","Number":4513},"contact":[{"role_name":"Owner","fname":"David","lname":"Lindberg"},{"role_name":"Owner","fname":"Rose","lname":"Lindberg"}]}

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
