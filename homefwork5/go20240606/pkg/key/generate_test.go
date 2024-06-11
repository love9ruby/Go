package key

import (
	"errors"
	"math/rand"
	"reflect"
	"testing"
)

func Test_generateKey(t *testing.T) {
	type args struct {
		randGen *rand.Rand
		result  chan string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test generateKey",
			args: args{
				randGen: rand.New(rand.NewSource(0)),
				result:  make(chan string, 1),
			},
			want: "mUNERA",
		},
		{
			name: "Test generateKey",
			args: args{
				randGen: rand.New(rand.NewSource(1)),
				result:  make(chan string, 1),
			},
			want: "BpLnfg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateKey(tt.args.randGen, tt.args.result)
		})
	}
}

func TestGenerateKeys(t *testing.T) {
	type args struct {
		num  int
		seed int64
		cb   func(s string) error
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test GenerateKeys",
			args: args{
				num:  5,
				seed: 0,
				cb: func(s string) error {
					if s == "mcjcEQ" {
						return errors.New("test")
					}
					return nil
				},
			},
			want: []string{"mUNERA", "9rI2cv", "TK4UHo", "vymkzA"},
		}, {
			name: "Test GenerateKeys",
			args: args{
				num:  -1,
				seed: 0,
				cb: func(s string) error {
					return nil
				},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateKeys(tt.args.num, tt.args.seed, tt.args.cb); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
