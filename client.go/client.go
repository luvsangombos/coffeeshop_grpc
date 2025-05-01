package main

import (
	pb "coffeeshop_grpc/coffee_shop_proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server ")
	}
	c := pb.NewCoffeeshopClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatal("error calling function getMenu")
	}

	done := make(chan bool)

	var items []*pb.Item

	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf(" can not receive %v", err)
			}

			items = resp.Item
			log.Printf("Resp received: %v", resp.Item)
		}
	}()
	<-done

	receipt, err := c.PlaceOrder(ctx, &pb.Order{Item: items})
	log.Printf("%v", receipt)
	status, err := c.GetOrderStatus(ctx, receipt)
	log.Printf("%v", status)
}
