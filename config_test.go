package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFileType(t *testing.T) {
	tests := []struct {
		name string
		c    string
		want string
	}{
		{name: "json", c: "./path_to/config.json", want: ".json"},
		{name: "yaml", c: "./path_to/config.yaml", want: ".yaml"},
		{name: "not support", c: "./path_to/config.txt", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkFileType(tt.c)
			if err != nil {
				assert.ErrorIs(t, err, errNotSupportFileType)
				return
			}
			if got != tt.want {
				t.Errorf("checkFileType() got = %v, want %v", got, tt.want)
			}
		})
	}
}
