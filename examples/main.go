package main

import (
	"encoding/json"
	"fmt"
	"log"

	apiv1 "github.com/mfridman/protoc-gen-go-json/examples/gen/go/api/v1"
)

func main() {
	by, err := json.Marshal(&apiv1.Request{
		Kind: &apiv1.Request_Name{
			Name: "alice",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(by))
	// {"name":"alice"}

	var request apiv1.Request
	if err := json.Unmarshal(by, &request); err != nil {
		log.Fatal(err)
	}
	fmt.Println(request.GetName())
	// alice
}
