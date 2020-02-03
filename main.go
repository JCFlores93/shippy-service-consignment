package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"log"
	"os"

	// Import the generated protobuf code
	pb "github.com/JCFlores93/shippy-service-consignment/proto/consignment"
	vesselProto "github.com/JCFlores93/shippy-service-vessel/proto/vessel"
)

const (
	port = ":50051"
	defaultHost = "datastore:27017"
)

func main() {
	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("shippy.service.consignment"),
		)

	// Init will parse the command line flags.
	srv.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	consignmentCollection := client.Database("shippy").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}
	vesselClient := vesselProto.NewVesselServiceClient("shippy.service.vessel", srv.Client())
	h := &handler{ repository, vesselClient}

	// Register handlers
	pb.RegisterShippingServiceHandler(srv.Server(), h)

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
