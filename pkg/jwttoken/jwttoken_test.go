package jwttoken_test

import (
	"net/http"
	"testing"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/jwttoken"
	"github.com/stretchr/testify/assert"
)

func TestAddJwtToken(t *testing.T) {
	type args struct {
		user string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Generate JWT token test...",
			args:    args{user: "serg"},
			wantErr: true,
		},
		{
			name:    "Generate JWT token test negative...",
			args:    args{user: ""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tokenStr, err := jwttoken.AddJwtToken(tt.args.user); (err == nil) != tt.wantErr {
				t.Errorf("AddJwtToken() error = %v, wantErr %v, token = %v", err, tt.wantErr, tokenStr)
			}
		})
	}
}

func TestValidateJwtToken(t *testing.T) {
	type args struct {
		user string
		in   *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Validate JWT token test...",
			args: args{in: &http.Request{Method: http.MethodPost,
				Header: http.Header{},
			},
				user: "TestUser",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Validate JWT token test negative...",
			args: args{in: &http.Request{Method: http.MethodPost,
				Header: http.Header{},
			},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenStr, _ := jwttoken.AddJwtToken(tt.args.user)
			tt.args.in.Header.Add("Jwt-Token", tokenStr)

			err, got := jwttoken.ValidateJwtToken(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJwtToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateJwtToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkAddJwtToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := jwttoken.AddJwtToken("TestUser")
		assert.NoError(b, err)
	}
}

func BenchmarkValidateJwtToken(b *testing.B) {
	tokenStr, _ := jwttoken.AddJwtToken("TestUser")
	for i := 0; i < b.N; i++ {
		err, _ := jwttoken.ValidateJwtToken(&http.Request{
			Method: http.MethodPost,
			Header: http.Header{
				"Jwt-Token": []string{tokenStr},
			}})
		assert.NoError(b, err)
	}
}
