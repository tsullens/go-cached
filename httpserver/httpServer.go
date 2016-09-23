package httpserver

import (
  "fugozi/database"
  "net/http"
  "sync"
  "log"
  "encoding/json"
  "time"
//  "fmt"
  "strings"
)

var (
  self *httpServer
  buckets = struct {
    sync.RWMutex
    m map[string]*database.Bucket
  }{m: make(map[string]*database.Bucket)}
)
const (
  timeLayout = "2006-01-02 15:04:05.00 MST"
)

type httpServer struct {
  IpAddr string
  Port string
  Status string
  StartTime string
  Debug bool
}

func NewHttpServer(args ...string) (*httpServer) {
  var ip, p string
  switch len(args){
  case 0:
    ip = ""
    p = ":3001"
  case 1:
    ip = args[0]
    p = ":3001"
  case 2:
    ip = args[0]
    p = args[1]
  }
  return &httpServer{
    IpAddr: ip,
    Port: p,
    Status: "Initialized",
    Debug: false,
  }
}

func (srv *httpServer) SetHttpServerDebug(val bool) {
  srv.Debug = val
}

func (srv *httpServer) RunServer() {

  srv.Status = "Running"
  srv.StartTime = time.Now().Format(timeLayout)
  self = srv

  initialize()

  // Route Handlers
  http.HandleFunc("/status/", statusHandler)
//  http.HandleFunc("/status/buckets/", bucketsHandler)
  http.HandleFunc("/bucket/", dbHandler)
  http.HandleFunc("/", rootHandler)
  log.Printf("Listening on %s", srv.Port)

  // Start the server
  listen := []string{srv.IpAddr, srv.Port}
  log.Fatal(http.ListenAndServe(strings.Join(listen, ""), nil))
}

// Route declarations
func rootHandler(w http.ResponseWriter, r *http.Request) {
  rlog("rootHandler", r)
  http.Redirect(w, r, "/status", http.StatusFound)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
  rlog("statusHandler", r)
  w.Header().Set("Content-Type", "application/json")
  js, err := json.MarshalIndent(&self, "", "  ")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Write(js)
}

func bucketsHandler(w http.ResponseWriter, r *http.Request) {
  rlog("bucketsHandler", r)
  w.Header().Set("Content-Type", "application/json")
  js, err := json.MarshalIndent(buckets, "", "  ")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Write(js)
}
