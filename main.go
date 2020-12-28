package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(context.Background())

	target := flag.Int("target", 0, "query running containers, add container till target reached")
	flag.Parse()
	fmt.Println(*target)
	for {
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if len(containers) < *target {
			fmt.Println("Trying to create container")
			_, err := cli.ContainerCreate(context.Background(), &container.Config{
				Tty:   true,
				Image: "ektor-client-scratch",
			}, &container.HostConfig{
				NetworkMode: "host",
				Resources: container.Resources{
					Memory: 30720000,
				},
			}, &network.NetworkingConfig{}, nil, "")

			if err != nil {
				fmt.Println("ERRORRR")
				fmt.Println(err)
			}
		}
		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			fmt.Printf("%s %s\n", container.ID[:10], container.Image)
		}
		time.Sleep(time.Minute)
	}
}
