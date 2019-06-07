package main

import (
  "log"
  "flag"
  "net/http"
  "strings"
)

// handler of static files
const prefixSlash = "/"
func FileHandler(w http.ResponseWriter, r *http.Request) {
  var path string
  var fileName string

  if strings.HasPrefix(r.URL.Path, prefixSlash) {
    path = r.URL.Path[len(prefixSlash):]
  }

  log.Println("The path is: ", path)
  if strings.HasPrefix(path, "ng") {
    fileName = "./" + path
  } else {
    fileName = "./ng/index.html"
  }
  log.Printf("Serveing files by FileHandler")
  
  http.ServeFile(w,r,fileName)
}
// \handler of static files


// msg type for websocket
type msg struct {
  Command string
}

// server entry
func main() {
  // handle request
  port := flag.String("p", "3000", "port")
  dir := flag.String("d", ".", "dir")
  flag.Parse()

  // serve static files
  http.HandleFunc(prefixSlash, FileHandler)
  
  log.Printf("Serving %s Http port: %s\n", *dir, *port)
  log.Fatal(http.ListenAndServe(":" + *port, nil))
}