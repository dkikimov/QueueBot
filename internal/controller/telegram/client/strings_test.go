package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"QueueBot/internal/entity"
)

func Test_cutStringByLines(t *testing.T) {
	type args struct {
		s           string
		linesToHave int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "OK",
			args: args{
				s: entity.ListToString([]entity.User{
					{ID: 1, Name: "User1"},
					{ID: 2, Name: "User2"},
					{ID: 3, Name: "User3"},
					{ID: 4, Name: "User4"},
					{ID: 5, Name: "User5"},
				}),
				linesToHave: 3,
			},
			want: "....\nUser3\nUser4\nUser5",
		},
		{
			name: "Lines to have is more than lines in string",
			args: args{
				s: entity.ListToString([]entity.User{
					{ID: 1, Name: "User1"},
					{ID: 2, Name: "User2"},
				}),
				linesToHave: 3,
			},
			want: entity.ListToString([]entity.User{
				{ID: 1, Name: "User1"},
				{ID: 2, Name: "User2"},
			}),
		},
		{
			name: "Lines to have is equal to lines in string",
			args: args{
				s: entity.ListToString([]entity.User{
					{ID: 1, Name: "User1"},
					{ID: 2, Name: "User2"},
				}),
				linesToHave: 2,
			},
			want: entity.ListToString([]entity.User{
				{ID: 1, Name: "User1"},
				{ID: 2, Name: "User2"},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, cutStringByLines(tt.args.s, tt.args.linesToHave), "cutStringByLinesWithCurrent(%v, %v)", tt.args.s, tt.args.linesToHave)
		})
	}
}

func Test_cutStringByLinesWithCurrent(t *testing.T) {
	type args struct {
		s                 string
		halfOfLinesToHave int
		currentIdx        int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "OK",
			args: args{
				s: entity.ListToString([]entity.User{
					{ID: 1, Name: "User1"},
					{ID: 2, Name: "User2"},
					{ID: 3, Name: "User3"},
					{ID: 4, Name: "User4"},
					{ID: 5, Name: "User5"},
				}),
				halfOfLinesToHave: 1,
				currentIdx:        2,
			},
			want: "....\nUser2\nUser3\nUser4\n....",
		},
		{
			name: "Lines to have is more than lines in string",
			args: args{
				s: entity.ListToString([]entity.User{
					{ID: 1, Name: "User1"},
					{ID: 2, Name: "User2"},
				}),
				halfOfLinesToHave: 2,
				currentIdx:        1,
			},
			want: entity.ListToString([]entity.User{
				{ID: 1, Name: "User1"},
				{ID: 2, Name: "User2"},
			}),
		},
		{
			name: "Lines to have is equal to lines in string",
			args: args{
				s: entity.ListToString([]entity.User{
					{ID: 1, Name: "User1"},
					{ID: 2, Name: "User2"},
				}),
				halfOfLinesToHave: 1,
				currentIdx:        1,
			},
			want: entity.ListToString([]entity.User{
				{ID: 1, Name: "User1"},
				{ID: 2, Name: "User2"},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				cutStringByLinesWithCurrent(tt.args.s, tt.args.halfOfLinesToHave, tt.args.currentIdx),
				"cutStringByLinesWithCurrent(%v, %v, %v)",
				tt.args.s,
				tt.args.halfOfLinesToHave,
				tt.args.currentIdx,
			)
		})
	}
}
