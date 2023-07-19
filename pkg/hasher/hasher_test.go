package hasher_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/hasher"
)

/* Sha256 from "TestPass" == eddef9e8e578c2a560c3187c4152c8b6f3f90c1dcf8c88b386ac1a9a96079c2c
 */
func TestCheckPasswordHash(t *testing.T) {
	want := "eddef9e8e578c2a560c3187c4152c8b6f3f90c1dcf8c88b386ac1a9a96079c2c"
	if !hasher.NewHasher("TestPass").CheckPasswordHash(want) {
		t.Errorf("CheckPasswordHash() => want %q", want)
	}

	if hasher.NewHasher("TestPass").CheckPasswordHash("any") {
		t.Errorf("CheckPasswordHash() => approved any hash)")
	}

}

/* Sha256 from "TestPass1" == 4ee33bac59675856c9d8f9ddfaf21368a08f8afe7827516c6d031b8859064229
 */
func TestHashPassword(t *testing.T) {
	want := "4ee33bac59675856c9d8f9ddfaf21368a08f8afe7827516c6d031b8859064229"
	if got, err := hasher.NewHasher("TestPass1").HashPassword(); got != want {
		fmt.Printf("%x", got)
		t.Errorf("HashPassword() = %q, want %q, err %s", got, want, err)
	}
}

func TestNewHasher(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name string
		args args
		want *hasher.HashingData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasher.NewHasher(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHasher() = %v, want %v", got, tt.want)
			}
		})
	}
}
