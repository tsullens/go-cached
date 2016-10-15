package server

import (
  "fugozi/database"
  "fugozi/util"
  "net/http"
  "sync"
  "encoding/json"
  "time"
  "fmt"
  "strings"
  "os"
)

var (
  self *httpServer
  buckets = struct {
    sync.RWMutex
    m map[string]*database.Bucket
  }{m: make(map[string]*database.Bucket)}
)
const (
  timeLayout = "2006-01-02 15:04:05.000 MST"
)

type httpServer struct {
  IpAddr string
  Port string
  Logger *util.LumberJack `json:"-"`
  Status string
  StartTime string
  Debug bool
}

func NewHttpServer(args ...string) (*httpServer) {
  var ip, p string
  lggr := util.NewLumberJack("db.log")
  switch len(args){
  case 0:
    ip = ""
    p = ":3341"
  case 1:
    ip = args[0]
    p = ":3341"
  case 2:
    ip = args[0]
    p = args[1]
  }
  return &httpServer{
    IpAddr: ip,
    Port: p,
    Logger: lggr,
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

  http.HandleFunc("/bucket/", dbHandler)
  http.HandleFunc("/", rootHandler)

  lgmsg := fmt.Sprintf("Listening on %s", srv.Port)
  self.Logger.Write(lgmsg)

  // Start the server
  listen := []string{srv.IpAddr, srv.Port}

  err := http.ListenAndServe(strings.Join(listen, ""), nil)
  self.Logger.Write(err.Error())
  os.Exit(1)
}

func RequestLog(msg string, start time.Time) {
  elapsed := time.Since(start)
  self.Logger.Write(fmt.Sprintf("%s %s", msg, elapsed))
}


// Route declarations
func rootHandler(w http.ResponseWriter, r *http.Request) {
//  rlog("rootHandler", r)
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  http.Redirect(w, r, "/status", http.StatusFound)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
//  rlog("statusHandler", r)
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  w.Header().Set("Content-Type", "application/json")
  js, err := json.MarshalIndent(&self, "", "  ")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  w.Write(js)
}

func bucketsHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  w.Header().Set("Content-Type", "application/json")
  js, err := json.MarshalIndent(buckets, "", "  ")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  w.Write(js)
}
