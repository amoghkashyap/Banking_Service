package main

import (
	"Alex_workshop/CassandraMigrations/cassandra"
	"Alex_workshop/CassandraMigrations/cassandra/entities"
	pb "Banking_Service/Interface"
	"context"
	"github.com/google/uuid"
	"strings"
	"log"
	"net"
	"google.golang.org/grpc"
)

const (
	port = ":17001"
)

var (
	deleteStatus bool
)

type server  struct {}

func (s *server ) CreateAccount(ctx context.Context, user *pb.CustomerDetails) (*pb.BankResponse, error) {
	person := entities.User{}
	person.SetName(user.GetName())
	person.SetAge(int(user.GetAge()))
	person.SetAddress(user.GetAddress())
	person.SetEmailID(user.GetEmailId())
	person.SetPassword(user.GetPassword())
	person.SetBalance(500)

	isAccountPresent := cassandra.FindAccountWithEmailId(person.GetEmailID())

	if isAccountPresent {
		return &pb.BankResponse{Response: "Account already exists, transaction ID: " + uuid.New().String(), Status: false}, nil
	} else {
		status := cassandra.InsertAccount(person)
		if status {
			return &pb.BankResponse{Response: "Account Created Successfully, transaction ID: " + uuid.New().String(), Status: true}, nil
		} else {
			return &pb.BankResponse{Response: "Account Creation failed , transaction ID: " + uuid.New().String(), Status: false}, nil
		}
	}
}

func (s *server) DeleteAccount(ctx context.Context, user *pb.CustomerDetails) (*pb.BankResponse, error) {
	isAccountPresent := cassandra.FindAccountWithEmailId(user.GetEmailId())

	if isAccountPresent {
		retreivePassword := cassandra.GetPassword(user.GetEmailId())
		if strings.EqualFold(retreivePassword, user.GetPassword()) {
			deleteStatus = cassandra.DeleteAccountWithEmailID(user.GetEmailId())
			if deleteStatus {
				return &pb.BankResponse{Response: "Account Closed Successfully, transaction ID: " + uuid.New().String(), Status: true}, nil
			} else {
				return &pb.BankResponse{Response: "Account Closure failed , transaction ID: " + uuid.New().String(), Status: false}, nil
			}
		} else {
			return &pb.BankResponse{Response: "Password did not match. Please try again, transaction ID: " + uuid.New().String(), Status: false}, nil
		}
	} else {
		return &pb.BankResponse{Response: "Account Does not exist. Please try again, transaction ID: " + uuid.New().String(), Status: false}, nil
	}
}

func (s *server) Deposit(ctx context.Context, user *pb.AuthenticationDetails) (*pb.BalanceResponse, error) {
	isAccountPresent := cassandra.FindAccountWithEmailId(user.GetEmailId())

	if isAccountPresent {
		retreivePassword := cassandra.GetPassword(user.GetEmailId())
		if strings.EqualFold(retreivePassword, user.GetPassword()) {
			retreiveBalance := cassandra.GetBalance(user.GetEmailId())
			newBalance := retreiveBalance + int(user.GetAmount())
			cassandra.UpdateBalance(newBalance,user.GetEmailId())
			return &pb.BalanceResponse{Balance: int32(newBalance), Response: "Amount successfully Deposited transaction ID: " + uuid.New().String()}, nil
		}
		return &pb.BalanceResponse{Response: "Password did not match. Please try again, transaction ID: " + uuid.New().String()}, nil
	} else {
		return &pb.BalanceResponse{Response: "Account Does not exist. Please try again, transaction ID: " + uuid.New().String()}, nil
	}
}

func (s *server) Withdraw(ctx context.Context, user *pb.AuthenticationDetails) (*pb.BalanceResponse, error) {
	isAccountPresent := cassandra.FindAccountWithEmailId(user.GetEmailId())

	if isAccountPresent {
		retreivePassword := cassandra.GetPassword(user.GetEmailId())
		if strings.EqualFold(retreivePassword, user.GetPassword()) {
			retreiveBalance := cassandra.GetBalance(user.GetEmailId())
			newBalance := retreiveBalance - int(user.GetAmount())
			if newBalance > 0 {
				cassandra.UpdateBalance(newBalance,user.GetEmailId())
				return &pb.BalanceResponse{Balance: int32(newBalance), Response: "Amount successfully Withdrawn transaction ID: " + uuid.New().String()}, nil
			} else {
				return &pb.BalanceResponse{Balance: int32(retreiveBalance), Response: "Insuffecient Funds transaction ID: " + uuid.New().String()}, nil
			}

		}
		return &pb.BalanceResponse{Response: "Password did not match. Please try again, transaction ID: " + uuid.New().String()}, nil
	} else {
		return &pb.BalanceResponse{Response: "Account Does not exist. Please try again, transaction ID: " + uuid.New().String()}, nil
	}
}

func main() {
	log.Println("starting Banking services")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBankServer(s,&server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal("error %v", err)
	}
}
