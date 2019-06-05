package apis

import (
  "log"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "bytes"
  "regexp"
)

// api
type defaultResp struct {
  Resp string
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    http.Error(w, "API only handles the requests via POST", http.StatusNotFound)
    return
  }

  matched, err := regexp.MatchString(`^/apis/(config|state)`, r.URL.Path)
  var respStr = "Unknown";
  if (matched) {
    respStr = "incomplete url /apis/config or apis/state without detail";
  }

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
  res := &defaultResp{Resp: respStr}

  js, err := json.Marshal(res)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}