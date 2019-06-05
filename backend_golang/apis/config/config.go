package config

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
type configResp struct {
  Key int
  Items []string
}
var BaseUrlOfConfigHandler = "/kube-apis/config/"

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
  // Testing of Redis operations
  name := "product"
  val := "scm"
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

  subUrl := strings.TrimPrefix(r.URL.Path, BaseUrlOfConfigHandler)

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
  res := &configResp{
    Key: 1,
    Items: []string{"config", subUrl, "three", string(b)}}

  js, err := json.Marshal(res)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}