package main

import (
	//"encoding/json"
	//"errors"
	"fmt"
	//"github.com/go-goodies/go_oops"
	//"github.com/l3x/jsoncfgo"
	"html/template"
	//"io/ioutil"
	"log"
	"net/http"
	//"regexp"
	"github.com/skycoin/skycoin/src/cipher"
	"os"
)

/*
func HtmlFileHandler(response http.ResponseWriter, request *http.Request, filename string) {
	response.Header().Set("Content-type", "text/html")
	webpage, err := ioutil.ReadFile(Dir + filename) // read whole the file
	if err != nil {
		http.Error(response, fmt.Sprintf("%s file error %v", filename, err), 500)
	}
	fmt.Fprint(response, string(webpage))
}
*/

//func IndexHandler(response http.ResponseWriter, request *http.Request) {
//	HtmlFileHandler(response, request, "/help.html")
//}

var server_pubkey cipher.PubKey = cipher.MustPubKeyFromHex("03f13f397a0a8d8840525c1b5c70eb21c8d49acb63cfa0d7225e57075ff4ec9e8a")
var server_seckey cipher.SecKey = cipher.MustSecKeyFromHex("6668f612441aa38da2ff4648ddbb6296ff91c126de7132126df88b3011510aa9")

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//t, _ := template.New("index.html").Parse(form)
	//t.Execute(w, "")

	t := template.Must(template.ParseFiles("index.html"))
	err := t.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	var text string
	if len(r.URL.Query()["TextBox1"]) == 0 {
		text = "empty"
	} else {
		text = r.URL.Query()["TextBox1"][0]
	}

	hash := cipher.SumSHA256([]byte(text))

	sig := cipher.SignHash(hash, server_seckey)
	text = fmt.Sprintf("%s\n===\nhash: %s\npubkey: %s\nsignature: %s\n", text, hash.Hex(), server_pubkey.Hex(), sig.Hex())

	data := map[string]string{"Text": text}

	fmt.Printf("===\n%v\n===\n", text)
	err = t.Execute(w, data)

}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("=== \n%v \n", r.URL.Query())

	fmt.Printf("TextBox1")

	in := r.URL.Query()["TextBox1"]
	_ = in
}

func main() {

	host := "localhost"
	fmt.Printf("host: %v\n", host)

	port := int(8080)
	fmt.Printf("port: %v\n", port)
	addr := fmt.Sprintf("%s:%d", host, port)

	mux := http.NewServeMux()

	mux.Handle("/event", http.HandlerFunc(eventHandler))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})

	err := http.ListenAndServe(addr, mux)
	fmt.Println(err.Error())
}
