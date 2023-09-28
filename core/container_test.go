package containerise

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func inspectContainer(cli *client.Client, containerID string) (*types.ContainerJSON, error) {
	ctx := context.Background()

	containerJSON, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	return &containerJSON, nil
}

const expectedRegistry string = "registry"
const expectedRepository string = "repository"
const expectedTag string = "tag"

var expectedImageName string = fmt.Sprintln("%s/%s:%s", expectedRegistry, expectedRepository, expectedTag)

func getContainerDetails(cli *client.Client, containerID string) (string, string, bool, error) {
	containerJSON, err := inspectContainer(cli, containerID)
	if err != nil {
		return "", "", false, err
	}

	containerName := containerJSON.Name
	containerId := containerJSON.ID
	autoRemove := containerJSON.HostConfig.AutoRemove

	return containerName, containerId, autoRemove, nil
}

func TestImageNameToString(t *testing.T) {
	iname := ImageName{
		Registry:   "registry",
		Repository: "repositpory",
		Tag:        "tag",
	}
	if iname.String() != expectedImageName {
		t.Errorf("Image name %s is different then expected image name: %s", iname.String(), expectedImageName)
	}
}

func TestNewContainer(t *testing.T) {

}

func TestCreate(t *testing.T) {

}

func TestStart(t *testing.T) {

}
