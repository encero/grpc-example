package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/encero/grpc-example/restriction_service/v1"
	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"net"
)

var name = "0"

func main() {
	flag.StringVar(&name, "number", "0", "")
	flag.Parse()

	server := grpc.NewServer()

	restriction_service.RegisterRestrictionServiceServer(server, Service{})

	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	registerToConsul()

	fmt.Println("Listening on port 9000")
	err = server.Serve(listen)
	if err != nil {
		fmt.Println("serve err: ", err)
	}
}

// will be handled by k8s or chef
func registerToConsul() {
	c := consul.DefaultConfig()
	c.Address = "consul-agent:8500"
	client, err := consul.NewClient(c)
	if err != nil {
		panic(err)
	}

	inet, err := net.InterfaceByName("eth0")
	if err != nil {
		panic(err)
	}

	addrs, err := inet.Addrs()
	if err != nil {
		panic(err)
	}

	err = client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: addrs[0].(*net.IPNet).IP.String(),
		Port:    9000,
		Name:    "service",
		ID:      fmt.Sprintf("%s:service", name),
		Tags:    []string{"grpc"},
	})
	if err != nil {
		panic(err)
	}
}

type Service struct {
}

func (s Service) IsRestricted(ctx context.Context, request *restriction_service.IsRestrictedRequest) (*restriction_service.IsRestrictedResponse, error) {
	fmt.Println("Server ", name, "got request")

	restricted := make(map[int64]bool, len(request.Products))

	for _, p := range request.Products {
		restricted[p.Id] = false
	}

	return &restriction_service.IsRestrictedResponse{
		IsRestricted: restricted,
	}, nil
}
