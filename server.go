package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "coffeeshop_grpc/coffee_shop_proto"
)

type server struct {
	pb.UnimplementedCoffeeshopServer
}

func (s *server) GetMenu(menuRequest *pb.MenuRequest, srv grpc.ServerStreamingServer[pb.Menu]) error {
	items := []*pb.Item{
		&pb.Item{
			Id:   "1",
			Name: "Black Coffee",
		},
		&pb.Item{
			Id:   "2",
			Name: "Americano",
		},
		&pb.Item{
			Id:   "3",
			Name: "Vanilla",
		},
	}

	for i, _ := range items {
		srv.Send(&pb.Menu{
			Item: items[0 : i+1],
		})
	}
	return nil

}

func (s *server) PlaceOrder(context context.Context, order *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "ABC",
	}, nil
}

func (s *server) GetOrderStatus(context context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status:  "IN_PROGRESS",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCoffeeshopServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %s", err)
	} 
	log.Println("Successfully running in port :9001")
	
}
