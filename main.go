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
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)

	target := flag.Int("target", 0, "query running containers, add container till target reached")
	flag.Parse()
	fmt.Println(*target)
	for {

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		numberToAdd := *target - len(containers)
		for i := 0; i < numberToAdd; i++ {
			fmt.Println("Trying to create container")
			resp, err := cli.ContainerCreate(ctx, &container.Config{
				Tty:   true,
				Image: "ektor-client-scratch",
			}, &container.HostConfig{
				NetworkMode: "host",
				Resources: container.Resources{
					Memory: 40720000,
				},
			}, &network.NetworkingConfig{}, nil, "")
			if err != nil {
				fmt.Println("ERRORRR")
				fmt.Println(err)
			}
			err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
			if err != nil {
				fmt.Println("ERRORRR")
				fmt.Println(err)
			} else {
				fmt.Println("Created")
			}
		}
		if err != nil {
			panic(err)
		}
		fmt.Println("Sleeping a minute")
		time.Sleep(time.Minute)
	}
}
