package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var (
	dEngine *DockerEngine
)

type DockerEngine struct {
	Cli *client.Client
}

type DockerCli interface {
	getAllDockersRunning() (string, error)
	createContainer() error
}

func initDockerEngine() error {
	newCli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	dEngine = &DockerEngine{Cli: newCli}
	return nil
}

func (d DockerEngine) getAllDockersRunning() (string, error) {
	containers, err := d.Cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	for _, container := range containers {
		line := fmt.Sprintf("%s %s\n", container.ID[:10], container.Image)
		sb.WriteString(line)
	}

	return sb.String(), nil
}

func (d DockerEngine) createContainer(req RequestCreateContainer) error {

	isImageFound, err := checkImageExist(d.Cli, req.ImageName)
	if err != nil {
		return err
	}
	if !isImageFound {
		errStr := fmt.Sprintf("image %s doesn't exist", req.ImageName)
		return errors.New(errStr)
	}

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8000",
	}
	containerPort, err := nat.NewPort("tcp", "80")
	if err != nil {
		return err
	}

	binding := fmt.Sprintf("%s:/home/laboratoire_user", req.FolderLocation)
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	containerName := getContainerName(req.ContainerName)
	fmt.Println(containerName)

	cont, err := d.Cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: req.ImageName,
		},
		&container.HostConfig{
			PortBindings: portBinding,
			Binds:        []string{binding},
		}, nil, containerName)
	if err != nil {
		return err
	}

	d.Cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})

	fmt.Printf("Container %s is started", cont.ID)
	return nil

}

func getContainerName(s string) string {
	if s == "" {
		t := time.Now()
		return t.Format("2006-01-0215-04-05")
	}
	return s
}

func checkImageExist(cli *client.Client, image string) (bool, error) {

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return false, err
	}

	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == image {
				// Image found locally
				fmt.Println("Image found locally:", image)
				fmt.Println("Image ID:", img.ID)
				return true, nil
			}
		}
	}
	return false, nil
}
