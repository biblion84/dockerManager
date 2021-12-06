package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var (
	target  = flag.Int("target", 10, "query running containers, add container till target reached")
	image   = flag.String("image", "", "Name of the image to be launched")
	verbose = flag.Bool("verbose", false, "Verbose")
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	ck(err)
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)

	flag.Parse()

	if *image == "" {
		flag.PrintDefaults()
		log.Fatal("Expected a docker image")
	}

	pt(target)
	for {

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		numberToAdd := *target - len(containers)
		for i := 0; i < numberToAdd; i++ {
			pt("Trying to create container")
			resp, err := cli.ContainerCreate(ctx, &container.Config{
				Tty:   true,
				Image: *image,
			}, &container.HostConfig{
				NetworkMode: "host",
				Resources: container.Resources{
					Memory: 50720000,
				},
			}, &network.NetworkingConfig{}, nil, "")

			ck(err)

			err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

			ck(err)

			if err == nil {
				pt("Created")
			}
		}
		ck(err)
		time.Sleep(time.Minute)
	}
}

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func pt(a ...interface{}) {
	if *verbose {
		fmt.Println(a)
	}
}
