package main

import (
	"google.golang.org/grpc"
	"log"
	pb "Banking_Service/Interface"
	"golang.org/x/net/context"
)

const (
	address = "localhost:17001"
)
func main() {
	conn, err := grpc.Dial(address,grpc.WithInsecure())
	if err != nil {
		log.Fatal(" error  %v", err)
	}
	defer conn.Close()
	client := pb.NewBankClient(conn)
	res, err := client.DeleteAccount(context.Background(),&pb.CustomerDetails{EmailId:"amogh.kashyap@nokia.com",Password:"Nokia"})
	if err != nil {
		log.Fatal("error  %v", err)
	}
	log.Println(res.Status)
	log.Println(res.Response)
}
