package appserver_test

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/appserver"
	"github.com/Serhii1Epam/enhanceHttpServer/pkg/userdata"
)

func TestAppserver_IsRun(t *testing.T) {
	tests := []struct {
		name string
		need bool
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "AppServer (false) test...",
			need: false,
			want: false,
		},
		{
			name: "AppServer (true) test...",
			need: true,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := appserver.NewServer()
			srv.Is_runned = tt.need
			if srv.IsRun() != tt.want {
				t.Errorf("IsRun error!")
			}
		})
	}
	/*srv := appserver.NewServer()
	if srv.IsRun() != false {

	}

	srv.Is_runned = true
	if srv.IsRun() != true {
		t.Errorf("IsRun (true) error!")
	}*/

}

func TestAppserver_isCorrectMethod(t *testing.T) {
	tests := []struct {
		name string
		in   *http.Request
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "POST method test...",
			in:   &http.Request{Method: http.MethodPost},
			want: true,
		},
		{
			name: "GET method test...",
			in:   &http.Request{Method: http.MethodGet},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := appserver.NewServer()
			if got := appserver.ExportIsCorrectMethod(srv, tt.in); got != tt.want {
				t.Errorf("Appserver.isCorrectMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppserver_userHandler(t *testing.T) {
	tests := []struct {
		name string
		in   *http.Request
		want error
	}{
		// TODO: Add test cases.
		{
			name: "Method GET...",
			in:   &http.Request{Method: http.MethodGet},
			want: errors.New("Server can't handle method [GET]. Continue...\n"),
		},
		{
			name: "Method POST...",
			in:   &http.Request{Method: http.MethodPost},
			want: errors.New("Cant parse incoming data"),
		},
		{
			name: "Method POST with JSON  Body...",
			in: &http.Request{Method: http.MethodPost,
				Body:   io.NopCloser(strings.NewReader("{\"user\": \"TestUser\", \"password\": \"TestPass\"}")),
				Header: http.Header{"Content-Type": {appserver.JSON}},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := appserver.NewServer()
			if tt.want == nil {
				if err := appserver.ExportUserHandler(srv, tt.in); err != nil {
					t.Errorf("Appserver.userHandler() error = %v, wantErr %v", err, tt.want)
				}
			} else {
				if err := appserver.ExportUserHandler(srv, tt.in); strings.Compare(err.Error(), tt.want.Error()) != 0 {
					t.Errorf("Appserver.userHandler() error = %v, wantErr %v", err, tt.want)
				}
			}
		})
	}
}

func TestAppserver_parseRequestBody(t *testing.T) {
	tests := []struct {
		name string
		in   *http.Request
		want *userdata.UserData
	}{
		// TODO: Add test cases.
		{
			name: "Method POST with TEXT Body...",
			in: &http.Request{Method: http.MethodPost,
				Body:   io.NopCloser(strings.NewReader("TestUser TestPass")),
				Header: http.Header{"Content-Type": {appserver.TEXT}},
			},
			want: &userdata.UserData{User: "TestUser", Password: "TestPass"},
		},
		{
			name: "Method POST with JSON Body witthout Header...",
			in: &http.Request{Method: http.MethodPost,
				Body: io.NopCloser(strings.NewReader("{\"user\": \"TestUser\", \"password\": \"TestPass\"}")),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := appserver.NewServer()
			if got := appserver.ExportParseRequestBody(srv, tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Appserver.parseRequestBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppserver_logMsgToConsole(t *testing.T) {
	type args struct {
		r *http.Request
		//w http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "POST text message test...",
			args: args{r: &http.Request{Method: http.MethodPost,
				Body:   io.NopCloser(strings.NewReader("TestUser TestPass")),
				Header: http.Header{"Content-Type": {appserver.TEXT}}},
			//w: httptest.NewRecorder(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := appserver.NewServer()
			appserver.ExportReqLogMsgToConsole(s, tt.args.r)
			//appserver.ExportRespLogMsgToConsole(s, tt.args.w)
		})
	}
}

func TestNewServer(t *testing.T) {
	srv := appserver.NewServer()
	if srv == nil {
		t.Errorf("appServer creation error!")
	}
}

/*func TestAppserver_logMsgToWriter(t *testing.T) {
	type args struct {
		r *http.Request
		w http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "POST text message test...",
			args: args{r: &http.Request{Method: http.MethodPost,
				Body:   io.NopCloser(strings.NewReader("TestUser TestPass")),
				Header: http.Header{"Content-Type": {appserver.TEXT}}},
				w: httptest.NewRecorder(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := appserver.NewServer()
			tt.args.w.WriteHeader(http.StatusOK)
			appserver.ExportLogMsgToWriter(s, tt.args.w, tt.args.r)
		})
	}
}*/
