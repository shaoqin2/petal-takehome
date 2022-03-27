package main

import "testing"

func Test_reverseString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "simple",
			s:    "abc",
			want: "cba",
		},
		{
			name: "with whitespace",
			s:    "abc def",
			want: "fed cba",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseString(tt.s); got != tt.want {
				t.Errorf("ReverseString() = %v, want %v", got, tt.want)
			}
		})
	}
}
