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
	res, err := client.Withdraw(context.Background(),&pb.AuthenticationDetails{EmailId:"amogh.kashyap@nokia.com",Password:"Nokia",Amount:1000})
	if err != nil {
		log.Fatal("error  %v", err)
	}
	log.Println(res.Balance)
	log.Println(res.Response)
}
