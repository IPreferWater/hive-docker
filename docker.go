package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var (
	cli client.Client
)

const BASE_FOLDER = "C:\\Users\\clems\\Workspace\\laboratoire"

func initCli() {
	newCli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	cli = *newCli
}

func getAllDockersRunning() string {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	res := ""
	for _, container := range containers {
		new := fmt.Sprintf("%s %s\n", container.ID[:10], container.Image)
		res = fmt.Sprintf("%s \n %s", res, new)
	}
	return res
}

func createContainer(req RequestCreateContainer) error {

	if !checkImageExist(req.ImageName) {
		errStr := fmt.Sprintf("image %s doesn't exist", req.ImageName)
		return errors.New(errStr)
	}

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8000",
	}
	containerPort, err := nat.NewPort("tcp", "80")
	if err != nil {
		panic("Unable to get the port")
	}

	binding := fmt.Sprintf("%s\\%s:/home/laboratoire_user", BASE_FOLDER, req.FolderLocation)
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: req.ImageName,
		},
		&container.HostConfig{
			PortBindings: portBinding,
			Binds:        []string{binding},
		}, nil, "name_for_container")
	if err != nil {
		panic(err)
	}

	cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})

	fmt.Printf("Container %s is started", cont.ID)
	return nil

}

//postgres:10.16
func checkImageExist(image string) bool {

	// List all images locally
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{
		//All:     false,
		//Filters: filters.NewArgs().ExactMatch("yadokari"),
	})
	if err != nil {
		panic(err)
	}

	// Iterate through the list of images and check if the provided image exists
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == image {
				// Image found locally
				fmt.Println("Image found locally:", image)
				fmt.Println("Image ID:", img.ID)
				return true
			}
		}
	}
	return false
}

func createVolume() {
	m := make(map[string]string)
	m["mount"] = "./test_folder"
	//m["destination"] = "/app"
	volumeBody := volume.VolumesCreateBody{
		//	Driver:     "local",

		Labels: m,
		Name:   "test_volume",
	}
	volume, err := cli.VolumeCreate(context.Background(), volumeBody)

	if err != nil {
		panic(err)
	}

	fmt.Printf("volume created %v", volume.Status)

}
