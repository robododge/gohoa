package gohoa

import "testing"

func Test_memberMapDB_AddConvience(t *testing.T) {
	type fields struct {
		mm map[Mkey]*Member
	}
	type args struct {
		member *Member
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"Add single member",
			fields{mm: make(map[Mkey]*Member)},
			args{member: &Member{PAddress: PropertyAddress{
				Number: 5566, StreetName: "Hawkins"},
				Contacts: []Contact{{RoleName: "Owner", FirstName: "Jan", LastName: "Dill"}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mdb := &memberMapDB{
			// 	mm: tt.fields.mm,
			// }

			mdb := NewMemberDB()
			mdb.(*memberMapDB).mm = tt.fields.mm

			mdb.AddConvience(tt.args.member)

			mOut, found := mdb.Fetch(MakeKey(tt.args.member))
			if !found {
				t.Error("Map boolean retruned not present entry")
			}
			if mOut == nil {
				t.Error("Memer is  Nil")
			}

			unknownMember := &Member{PAddress: PropertyAddress{Number: 5566, StreetName: "Jilion"}}
			mNot, found := mdb.Fetch(MakeKey(unknownMember))
			if mNot != nil || found {
				t.Error("Unknown member was not suppose to be present")
			}
		})
	}

}
func Test_AddingMutipleMemebrs(t *testing.T) {
	type args struct {
		memberDb MemberDB
		members  []Member
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test adding multiple members",
			args: args{
				memberDb: NewMemberDB(),
				members: []Member{
					{PAddress: PropertyAddress{
						Number: 5566, StreetName: "Westlane"},
						Contacts: []Contact{{RoleName: "Owner", FirstName: "Jim", LastName: "Dill"}}},
					{PAddress: PropertyAddress{
						Number: 4567, StreetName: "Hawkins"},
						Contacts: []Contact{{RoleName: "Owner", FirstName: "Jan", LastName: "Hill"}}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, m := range tt.args.members {
				locMember := m
				//Type assertion test
				_, ok := tt.args.memberDb.(*memberMapDB)
				if !ok {
					t.Errorf("memberMapDB type was expected but not found, type: %T", tt.args.memberDb)
				}
				tt.args.memberDb.AddConvience(&locMember)

			}

			expectedSize := len(tt.args.members)
			if tt.args.memberDb.Size() != expectedSize {
				t.Errorf(" Expected size is %d but found %d\n", expectedSize, tt.args.memberDb.Size())
			}

			if _, found := tt.args.memberDb.Fetch(MakeKeyV("Westlane", 5566)); !found {
				t.Error("Could not find expected member Westlane")
			}
			if _, found := tt.args.memberDb.Fetch(MakeKeyV("Hawkins", 4567)); !found {
				t.Error("Could not find expected member Hawkins")
			}

			w, _ := tt.args.memberDb.Fetch(MakeKeyV("Westlane", 5566))
			h, _ := tt.args.memberDb.Fetch(MakeKeyV("Hawkins", 4567))

			if w.Contacts[0].FirstName == h.Contacts[0].FirstName {
				t.Error("both contacts should not be the same, probably re-using a variable for pointer")
			}

		})

	}
}
