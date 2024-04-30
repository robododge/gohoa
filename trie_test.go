package gohoa

import (
	"reflect"
	"testing"
)

func TestTrie_TextSuggestions(t *testing.T) {
	trie := InitTrie(TEXT_TRIE)
	trie.Insert("apple")
	trie.Insert("banana")
	trie.Insert("orange")
	trie.Insert("mackey")

	trie.Insert("maclace")
	trie.Insert("macbig")
	trie.Insert("muadib")
	trie.Insert("mallory")

	tests := []struct {
		name       string
		key        string
		wantResult []string
	}{
		{
			name:       "Test suggestions for 'a'",
			key:        "a",
			wantResult: []string{"apple"},
		},
		{
			name:       "Test suggestions for 'b'",
			key:        "b",
			wantResult: []string{"banana"},
		},
		{
			name:       "Test suggestions for 'o'",
			key:        "o",
			wantResult: []string{"orange"},
		},
		{
			name:       "Test suggestions for 'ap'",
			key:        "ap",
			wantResult: []string{"apple"},
		},
		{
			name:       "Test suggestions for 'ban'",
			key:        "ban",
			wantResult: []string{"banana"},
		},
		{
			name:       "Test suggestions for 'ora'",
			key:        "ora",
			wantResult: []string{"orange"},
		},
		{
			name:       "Test suggestions for 'xyz'",
			key:        "xyz",
			wantResult: nil,
		},
		{
			name:       "Test suggestions for 'mac'",
			key:        "mac",
			wantResult: []string{"macbig", "mackey", "maclace"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := trie.Suggestions(tt.key)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Trie.Suggestions() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestTrie_NumberSuggestions(t *testing.T) {
	trie := InitTrie(NUM_TRIE)
	trie.Insert("123")
	trie.Insert("456")
	trie.Insert("789")
	trie.Insert("5698")
	trie.Insert("5611")
	trie.Insert("5612")

	tests := []struct {
		name       string
		key        string
		wantResult []string
	}{
		{
			name:       "Test suggestions for '1'",
			key:        "1",
			wantResult: []string{"123"},
		},
		{
			name:       "Test suggestions for '4'",
			key:        "4",
			wantResult: []string{"456"},
		},
		{
			name:       "Test suggestions for '7'",
			key:        "7",
			wantResult: []string{"789"},
		},
		{
			name:       "Test suggestions for '12'",
			key:        "12",
			wantResult: []string{"123"},
		},
		{
			name:       "Test suggestions for '45'",
			key:        "45",
			wantResult: []string{"456"},
		},
		{
			name:       "Test suggestions for '78'",
			key:        "78",
			wantResult: []string{"789"},
		},
		{
			name:       "Test suggestions for '999'",
			key:        "999",
			wantResult: nil,
		},
		{
			name:       "Test suggestions for '56'",
			key:        "56",
			wantResult: []string{"5611", "5612", "5698"},
		},
		{
			name:       "Test suggestions for '561'",
			key:        "561",
			wantResult: []string{"5611", "5612"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := trie.Suggestions(tt.key)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Trie.Suggestions() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
