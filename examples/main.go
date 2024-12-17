package main

import (
	"encoding/json"
	"fmt"
	"log"

	apiv1 "github.com/mfridman/protoc-gen-go-json/examples/gen/go/api/v1"
	"google.golang.org/protobuf/proto"
)

func main() {
	// In the opaque version, the message is not directly accessible. See
	// https://protobuf.dev/reference/go/opaque-faq/#builders-vs-setters for more information.
	//
	// 1. Use builder to create a request
	req1 := apiv1.Request_builder{
		Name: proto.String("alice"),
	}.Build()
	// 2. Use setter to create a request
	// req2 := &apiv1.Request{}
	// req2.SetName("alice")
	by1, err := json.Marshal(req1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("opaque:", string(by1))
	// opaque: {"name":"alice"}

	// In the open version, the message is directly accessible.
	by2, err := json.Marshal(&apiv1.Request{
		Kind: &apiv1.Request_Name{
			Name: "alice",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("open:", string(by2))
	// open: {"name":"alice"}

	var request apiv1.Request
	if err := json.Unmarshal(by2, &request); err != nil {
		log.Fatal(err)
	}
	fmt.Println(request.GetName())
	// alice
}
