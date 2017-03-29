package main

import (
  "net/http"
	"log"
	"github.com/gorilla/mux"
  "encoding/json"
  "io/ioutil"
  "gopkg.in/mgo.v2"
  "time"
  "fmt"
)

type Message struct{
  Name  string `json:"name",bson:"name",binding:"required"`
  Email string `json:"email",bson:"email",binding:"required"`
  Website string `json:"website",bson:"website",binding:"required"`
  Text  string `json:"text",bson:"text",binding:"required"`
  CreateDate time.Time  `json:"createDate",bson:"createDate",binding:"required"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func testjson(w http.ResponseWriter, r *http.Request) {
  b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
  log.Println("GET//JSON")
  w.Write(b)
}

func SendMessage(w http.ResponseWriter,r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  session, err := mgo.Dial("localhost")
  if err != nil {
      panic(err)
  }
  session.SetMode(mgo.Monotonic, true)
  MS := session.DB("webcontact").C("message")
  var m Message
  b, _ := ioutil.ReadAll(r.Body)
  json.Unmarshal(b, &m)
  m.CreateDate=time.Now()
  log.Println(m)
  err = MS.Insert(&m)
  if err!=nil{
    panic(err)
  }
  defer session.Close()
  w.Write(b)
}

func GetMessage(w http.ResponseWriter,r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  session, err := mgo.Dial("localhost")
  if err != nil {
      panic(err)
  }
  session.SetMode(mgo.Monotonic, true)
  MS := session.DB("webcontact").C("message")
  var m []Message
  err = MS.Find(nil).All(&m)
  if err!=nil{
    panic(err)
  }
  defer session.Close()
  log.Println(m)
  fmt.Fprint(w,m)
  w.Write(m)
}

func main() {
  // session, err := mgo.Dial("mongodb://Admin:pa22word@ds139979.mlab.com:39979/decive")

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", Handler).Methods("GET")
  r.HandleFunc("/json",testjson).Methods("GET")
  r.HandleFunc("/contact/getms",GetMessage).Methods("GET")
  r.HandleFunc("/contact/sendms",SendMessage).Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
