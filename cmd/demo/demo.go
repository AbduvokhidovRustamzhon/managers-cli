package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type notes struct {
	To string `xml:"to"`
	From string `xml:"from"`
	Heading string `xml:"heading"`
	Body string `xml:"body"`
}


type notew struct {
	To string `json:"to"`
	From string `json:"from"`
	Heading string `json:"heading"`
	Body string `json:"body"`
}


// Так можно делать в тесте, если вы передаёте io.Reader, io.Writer
func main() {
	file, err := os.Open("commands.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var login string
	fmt.Fscan(file, &login)
	fmt.Println(login)
	var password string
	fmt.Fscan(file, &password)
	fmt.Println(password)

note := &notes{
	To:      "Nicky",
	From:    "Rock",
	Heading: "Meeting",
	Body:    "Meeting at 5pm",
}

files,_ := xml.MarshalIndent(note,""," ")
_=ioutil.WriteFile("users.xml",files,0644)


	notek := &notew{
		To:      "Nicky",
		From:    "Rock",
		Heading: "Meeting",
		Body:    "Meeting at 5pm",
	}

	filew,_ := json.MarshalIndent(notek,""," ")
	_=ioutil.WriteFile("users.json",filew,0644)



}
