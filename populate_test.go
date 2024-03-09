package gohoa

import (
	"testing"
)

func TestAllMembers_deDupeMembers(t *testing.T) {
	// Create a sample AllMembers instance with duplicate members
	am := &AllMembers{
		Members: []Member{
			{MemberId: 771, MemberName: "John K 1", PAddress: PropertyAddress{StreetName: "Main St", Number: 123}},
			{MemberId: 200, MemberName: "John K 2", PAddress: PropertyAddress{StreetName: "Main St", Number: 123}},
			{MemberId: 1, MemberName: "John JJ", PAddress: PropertyAddress{StreetName: "Main St", Number: 777}},
			{MemberId: 3, MemberName: "Mike P", PAddress: PropertyAddress{StreetName: "Elm St", Number: 456}},
		},
	}

	// Call the deDupeMembers method
	am.DeDupeMembers()

	// Verify the expected number of members after deduplication
	expectedSize := 3
	if len(am.Members) != expectedSize {
		t.Errorf("Expected %d members after deduplication, but got %d", expectedSize, len(am.Members))
	}

	// Verify that the duplicate member with ID "1" is removed
	for _, m := range am.Members {
		if m.MemberId == 200 {
			t.Errorf("Duplicate member with ID '1' should be removed, but found member: %s", m.MemberName)
		}
	}

	// Verify that the member with the highest ID is kept
	highestID := int32(771)
	foundHighestID := false
	for _, m := range am.Members {
		if m.MemberId == highestID {
			foundHighestID = true
			break
		}
	}
	if !foundHighestID {
		t.Errorf("Member with highest ID '%d' should be kept, but not found", highestID)
	}

	// // Verify that the log messages are printed for duplicate members
	// expectedLogMessages := []string{
	// 	"OOPS! -Duplicate member John id:1",
	// 	" - Keeping member John id:1",
	// }
	// for _, msg := range expectedLogMessages {
	// 	if !logContains(msg) {
	// 		t.Errorf("Expected log message '%s' not found", msg)
	// 	}
	// }
}

// Helper function to check if a log message is printed
// func logContains(msg string) bool {
// 	// Implement your logic to check if the log contains the given message
// 	// You can use a log hook or redirect log output to a buffer for testing
// 	return false
// }
