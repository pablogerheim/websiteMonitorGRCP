package main

import (
	"fmt"
	"log"
	"net"
	"websiteMonitor/config"
	"websiteMonitor/pb"
	"websiteMonitor/resource"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	resource.Inject()
	fmt.Println("tcp", ":"+config.PORT)

	lis, err := net.Listen("tcp", ":"+config.PORT)
	if err != nil {
		log.Fatalf("Failed to listen on %v: %v\n", config.PORT, err)
	}

	s := grpc.NewServer()

	reflection.Register(s)

	pb.RegisterWebsiteMonitorServiceServer(s, resource.ServersControllers.Site)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
