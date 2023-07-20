package wsserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/wsserver"
)

func TestWsData_StartWs(t *testing.T) {
	type args struct {
		out http.ResponseWriter
		in  *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Start WebSokets...",
			args: args{
				in: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"Connection":            []string{"Upgrade"},
						"Upgrade":               []string{"websocket"},
						"Sec-Websocket-Key":     []string{"dGhlIHNhbXBsZSBub25jZQ=="},
						"Sec-Websocket-Version": []string{"13"},
					},
				},
				out: httptest.NewRecorder(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd := wsserver.NewWsData()
			if err := wd.StartWs(tt.args.out, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("WsData.StartWs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWsData_openWs(t *testing.T) {
	type args struct {
		out http.ResponseWriter
		in  *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Open WebSokets...",
			args: args{
				in: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"Connection":            []string{"Upgrade"},
						"Upgrade":               []string{"websocket"},
						"Sec-Websocket-Key":     []string{"dGhlIHNhbXBsZSBub25jZQ=="},
						"Sec-Websocket-Version": []string{"13"},
					},
				},
				out: httptest.NewRecorder(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd := wsserver.NewWsData()
			wd.StartWs(tt.args.out, tt.args.in)
			if err := wd.OpenWs(); (err != nil) != tt.wantErr {
				t.Errorf("WsData.openWs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
