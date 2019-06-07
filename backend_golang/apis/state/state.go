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
  fmt.Println("name (from Redis slave): ", respVal)
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

  // Parse request body
  var data = make(map[string]interface{})
  json.Unmarshal([]byte(b), &data)

  // a string slice to hold the keys
  var keys = make([]string, len(data))

  // iteration counter
  i := 0

  // copy data's keys into keys
  for s, _ := range data {
    keys[i] = s
    i++
  }

  // output result to STDOUT
  fmt.Printf("%#v\n", keys)
  // \Parse request body

  // print body string
  rbody := ioutil.NopCloser(bytes.NewBuffer(b))
  log.Printf("BODY: %q", rbody)


  // reply to the request
  res := &stateResp{
    Key: 1,
    Items: []string{"state", subUrl, string(b)}}

  js, err := json.Marshal(res)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}