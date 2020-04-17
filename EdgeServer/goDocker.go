package main

import (
	"context"
	"fmt"
//	"io"
//	"os"

	"github.com/docker/docker/api/types"
        "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
       "github.com/docker/go-connections/nat"
)
func stopContainer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}
}
func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	imageName := "tnreddy9/camera_client"

/*	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
*/
        hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8179",
	}
	containerPort := nat.Port("8180/tcp")
	if err != nil {
		panic("Unable to get the port")
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
//	portBinding := nat.PortMap{containerPort: hostBinding}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, &container.HostConfig{
                PortBindings: portBinding,
                Mounts: []mount.Mount{
                        { 
                            Type: mount.TypeBind,
                            Source: "/tmp/data",
                            Target: "/tmp/data",
                        },
              },
        }, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)
}
