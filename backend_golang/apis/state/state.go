package state

import (
  "log"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "bytes"
  "strings"
  "fmt"
  "dbs/rdHandler"
)

// api response
type stateResp struct {
  Key int
  Items []string
}
var BaseUrlOfStateHandler = "/kube-apis/state/"


func StateHandler(w http.ResponseWriter, r *http.Request) {
  // Testing of Redis operations
  name := "name"
  val := "alan"
  redisWrite := rdHandler.RedisdbWrite {}
  redisRead := rdHandler.RedisdbRead {}

  errWrite := redisWrite.Set(&name, &val);
  if errWrite != nil {
    panic(errWrite)
  }
  respVal, errRead := redisRead.Get(&name);
  if errRead != nil {
    panic(errRead)
  }
  fmt.Println("name: ", respVal)
  // \Testing of Redis operations

  if r.Method != "POST" {
    http.Error(w, "API only handles the requests via POST", http.StatusNotFound)
    return
  }

  subUrl := strings.TrimPrefix(r.URL.Path, BaseUrlOfStateHandler)

  // Read body
  b, err := ioutil.ReadAll(r.Body)
  defer r.Body.Close()
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }
  rbody := ioutil.NopCloser(bytes.NewBuffer(b))
  log.Printf("BODY: %q", rbody)

  // reply to the request
  res := &stateResp{
    Key: 1,
    Items: []string{"state", subUrl, "three", string(b)}}

  js, err := json.Marshal(res)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}