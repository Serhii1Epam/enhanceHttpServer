// Server package
// simple HTTP server for accepting user request
// check users passwords
package appserver

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/appdb"
	"github.com/Serhii1Epam/enhanceHttpServer/pkg/hub"
	"github.com/Serhii1Epam/enhanceHttpServer/pkg/jwttoken"
	"github.com/Serhii1Epam/enhanceHttpServer/pkg/userdata"
	"github.com/Serhii1Epam/enhanceHttpServer/pkg/wsserver"
)

type responseStatus struct {
	code int
	text string
}

func NewResponseError() *responseStatus {
	return &responseStatus{}
}

type Appserver struct {
	mux       *http.ServeMux
	db        *appdb.Database
	user      *userdata.UserData
	respStat  *responseStatus
	wsd       *wsserver.WsData
	wsHub     *hub.Hub
	Is_runned bool
}

// Content-Type constants
const (
	JSON = "application/json"
	TEXT = "text/plain"
)

// For testing
var (
	ExportIsCorrectMethod     = (*Appserver).isCorrectMethod
	ExportParseRequestBody    = (*Appserver).parseRequestBody
	ExportUserHandler         = (*Appserver).userHandler
	ExportLogMsgToWriter      = (*Appserver).logMsgToWriter
	ExportReqLogMsgToConsole  = (*Appserver).logReqMsgToConsole
	ExportRespLogMsgToConsole = (*Appserver).logRespMsgToConsole
)

//End for testing

func NewServer() *Appserver {
	return &Appserver{Is_runned: false}
}

func (s Appserver) IsRun() bool {
	return s.Is_runned
}

func (s Appserver) writeResponseError(w http.ResponseWriter, text string, code int) {
	s.setResponseStatus(code, text)
	http.Error(w, s.respStat.text, s.respStat.code)
}

func (s Appserver) setResponseStatus(code int, text string) {
	s.respStat.code = code
	s.respStat.text = text
}

func (s *Appserver) SrvRun() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Printf("$PORT is not set, use default...")
		port = "8080"
	}
	s.Is_runned = true
	s.db = appdb.NewDatabase()
	s.respStat = NewResponseError()
	s.mux = http.NewServeMux()
	s.wsd = wsserver.NewWsData()
	s.wsHub = hub.NewHub()
	go s.wsHub.StartHub()

	userLoginHandler := func(w http.ResponseWriter, r *http.Request) {
		if err := s.userHandler(r); err != nil {
			s.writeResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.user.Login(s.db); err != nil {
			s.writeResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//Add JWT token to headers
		tokenStr, err := jwttoken.AddJwtToken(s.user.User)
		if err != nil {
			s.writeResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Jwt-Token", tokenStr)
		//Redirect to chat WebSocket
		//s.setResponseStatus(http.StatusMovedPermanently, http.StatusText(http.StatusMovedPermanently))
		//http.Redirect(w, r, "/ws", s.respStat.code)
		//For testing return OK 200
		s.setResponseStatus(http.StatusSwitchingProtocols, http.StatusText(http.StatusSwitchingProtocols))
		w.WriteHeader(s.respStat.code)
	}

	userHandler := func(w http.ResponseWriter, r *http.Request) {
		if err := s.userHandler(r); err != nil {
			s.writeResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.user.Create(s.db); err != nil {
			s.writeResponseError(w, err.Error(), http.StatusInternalServerError)
		}
		s.setResponseStatus(http.StatusOK, http.StatusText(http.StatusOK))
		w.WriteHeader(s.respStat.code)
	}

	aboutHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple HTTP Server developed for GO switch program.\n")
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		s.writeResponseError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	wsHandler := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Upgrade connection to WebSocket... ")
		s.wsd.StartWs(s.wsHub, w, r)
	}

	s.mux.Handle("/about", s.middelwareLog(http.HandlerFunc(aboutHandler)))
	s.mux.Handle("/user/login", s.middelwareLog(http.HandlerFunc(userLoginHandler)))
	s.mux.Handle("/user", s.middelwareLog(http.HandlerFunc(userHandler)))
	s.mux.Handle("/ws", s.middelwareLog(s.authHandler(http.HandlerFunc(wsHandler))))
	s.mux.Handle("/", s.middelwareLog(http.HandlerFunc(indexHandler)))

	log.Fatal(http.ListenAndServe(":"+port, s.mux))
}

func (s *Appserver) logMsgToWriter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL [%s] ", r.URL.Path)
	fmt.Fprintf(w, "Method [%v] ", r.Method)
	fmt.Fprintf(w, "Content-Type: \"%s\"\n", r.Header.Get("Content-Type"))
}

func (s *Appserver) logReqMsgToConsole(r *http.Request) {
	log.Printf("--> %s %s %s", r.RemoteAddr, r.Method, r.URL)
}

func (s *Appserver) logRespMsgToConsole(w http.ResponseWriter) {
	log.Printf("<-- %s %d", s.respStat.text, s.respStat.code)
}

func (s *Appserver) middelwareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logReqMsgToConsole(r)
		defer func() {
			s.logMsgToWriter(w, r)
			s.logRespMsgToConsole(w)
		}()
		next.ServeHTTP(w, r)

	})
}

func (s *Appserver) authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Auth middelware... ")
		if err, ok := jwttoken.ValidateJwtToken(r); !ok {
			log.Printf("FAILED.\n")
			s.writeResponseError(w, err.Error(), http.StatusUnauthorized)
			return
		} else {
			log.Printf("OK.\n")
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Appserver) parseRequestBody(r *http.Request) *userdata.UserData {
	var str userdata.IParse
	if r.Body == nil {
		return nil
	}
	req, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	r.Body.Close()
	switch r.Header.Get("Content-Type") {
	case JSON:
		{
			str = userdata.JsonBytes(req)

		}
	case TEXT:
		{
			str = userdata.PlainTextBytes(req)
		}
	default:
		{
			return nil
		}
	}
	return userdata.Parse(str)
}

func (s *Appserver) isCorrectMethod(r *http.Request) bool {
	switch r.Method {
	case http.MethodPost:
		{
			return true
		}
	default:
		{
			return false
		}
	}
}

func (s *Appserver) userHandler(r *http.Request) error {
	if !s.isCorrectMethod(r) {
		customMsg := fmt.Sprintf("Server can't handle method [%v]. Continue...\n", r.Method)
		return errors.New(customMsg)
	}
	if s.user = s.parseRequestBody(r); s.user == nil {
		return errors.New("Cant parse incoming data")
	}

	return nil
}
