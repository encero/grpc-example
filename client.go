package main

import (
	"context"
	"fmt"
	"github.com/encero/grpc-example/restriction_service/v1"
	"google.golang.org/grpc"
	"os"
	"time"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

func main() {
	connectionString := os.Getenv("GRPC_CONNECTION_STRING")
	if connectionString == "" {
		connectionString = "consul://consul-agent:8500/service?healthy=true"
	}

	conn, err := grpc.Dial(
		connectionString,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		panic(err)
	}

	client := restriction_service.NewRestrictionServiceClient(conn)

	for i := 0; ; i++ {
		start := time.Now()
		resp, err := client.IsRestricted(context.Background(), &restriction_service.IsRestrictedRequest{
			Products: []*restriction_service.Product{
				{
					Id:       1000111011,
					Brand:    "INTEX",
					Category: "ND001",
				},
			},
		})
		took := time.Since(start)
		if err != nil {
			panic(err)
		}

		for id, result := range resp.IsRestricted {
			fmt.Println(id, "is restricted", result)
		}

		fmt.Printf("got response: roundtrip took :%d us\n", took.Microseconds())
		time.Sleep(time.Millisecond * 100)
	}

}
