package rules

import "testing"

func TestPostFilter(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// Failure cases
		{"Empty String", args{""}, false},
		{"Starts with lowercase", args{"hello"}, false},
		{"Contains emoji", args{"Hello 😀."}, false},
		{"Contains url", args{"Hello http://example.com."}, false},
		{"Contains url", args{"Hello www.example."}, false},
		{"Ends with ellipsis", args{"Hello..."}, false},
		{"Contains hashtag", args{"Hello #world."}, false},
		{"Contains @mention", args{"Hello @world."}, false},
		{"Multiple repeated ?", args{"Hello??."}, false},
		{"Multiple repeated !", args{"Hello!!."}, false},
		{"NO YELLING", args{"HELLO."}, false},
		// Success cases
		{"Starts with capital letter and ends with period", args{"Hello."}, true},
		{"Ends with question mark", args{"Hello?"}, true},
		{"Ends with exclamation mark", args{"Hello!"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PostFilter(tt.args.text); got != tt.want {
				t.Errorf("PostFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
