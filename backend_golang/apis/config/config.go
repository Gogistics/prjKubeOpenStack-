package config

import (
  // "log"
  // "bytes"
  "net/http"
  "encoding/json"
  "io/ioutil"
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
  if r.Method == "POST" {
    post(w, r)
  } else if r.Method == "GET" {
    get(w, r)
  } else {
    http.Error(w, "API only handles the requests via POST", http.StatusNotFound)
    return
  }
}


func post(w http.ResponseWriter, r *http.Request) {
  redisWrite := rdHandler.RedisdbWrite {}
  hash := strings.TrimPrefix(r.URL.Path, BaseUrlOfConfigHandler)

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
  for key, _ := range data {
    keys[i] = key
    i++

    // write data into redis secondary index
    // fmt.Printf("data[key]: %q", data[key])
    val := fmt.Sprintf("%v", data[key])
    redisWrite.HSet(&hash, &key, &val);
  }

  // output result to STDOUT
  // fmt.Printf("%#v\n", keys)
  // \Parse request body

  // print body string
  // rbody := ioutil.NopCloser(bytes.NewBuffer(b))
  // log.Printf("BODY: %q", rbody)

  // reply to the request
  res := &configResp{
    Key: 1,
    Items: []string{"config", hash, string(b)}}

  js, err := json.Marshal(res)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func get(w http.ResponseWriter, r *http.Request) {
  redisRead := rdHandler.RedisdbRead {}
  // read data
  hash := strings.TrimPrefix(r.URL.Path, BaseUrlOfConfigHandler)
  // get cf3 status
  respMap, errRead := redisRead.HGetAll(&hash);
  if errRead != nil {
    panic(errRead)
  }
  fmt.Println("states: ", respMap)
  js, err := json.Marshal(respMap)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}