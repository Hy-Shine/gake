package main

import (
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		e    string
		want bool
	}{
		{name: "contains", s: []string{"a", "b", "c"}, e: "b", want: true},
		{name: "not contains", s: []string{"a", "b", "c"}, e: "d", want: false},
		{name: "empty", s: []string{}, e: "d", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.s, tt.e); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplication(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		want []string
	}{
		{name: "remove duplication", s: []string{"a", "b", "c", "a", "b", "c"}, want: []string{"a", "b", "c"}},
		{name: "empty", s: []string{}, want: []string{}},
		{name: "only one", s: []string{"a"}, want: []string{"a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := distinct(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDuplication() = %v, want %v", got, tt.want)
			}
		})
	}
}
