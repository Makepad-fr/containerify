package containerise

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var cli *client.Client

func init() {
	// Create a new Docker client
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli = c
}

type ContainerStatus struct {
	Started bool
	Created bool
	Removed bool
}

type ContainerConfig struct {
	StopTimeout *int
	AutoRemove  bool
	User        string
}

type ImageName struct {
	Registry   string
	Repository string
	Tag        string
}

func (i ImageName) String() string {
	return fmt.Sprintf("%s/%s:%s", i.Registry, i.Repository, i.Tag)
}

type Container struct {
	ID        *string
	name      string
	imageName ImageName
	ctx       context.Context
	config    ContainerConfig
	status    ContainerStatus
}

// NewContainer initialise a new *Container instance with given ImagEName, container name and ContainerConfig.
// It returns an error if seomthing foes wrong when pulling the image
func NewContainer(ctx context.Context, imageName ImageName, name string, config ContainerConfig) (*Container, error) {
	// Pull the image (optional, if the image is not already present)
	_, err := cli.ImagePull(ctx, imageName.String(), types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	var container Container = Container{
		ID:        nil,
		name:      name,
		imageName: imageName,
		ctx:       ctx,
		config:    config,
		status: ContainerStatus{
			Started: false,
			Created: false,
			Removed: false,
		},
	}
	return &container, nil
}

// Create creates a new container with the ImageName and the ContainerConfig
// it returns an error if somethings goes wrong when creating the container
func (s *Container) Create(command []string) error {
	// Create a new container with the --rm option
	containerCreateResponse, err := cli.ContainerCreate(
		s.ctx,
		&container.Config{
			Image: s.imageName.String(),
			Cmd:   command,
			User:  s.config.User,
		},
		&container.HostConfig{
			AutoRemove: s.config.AutoRemove,
		},
		nil,
		nil,
		s.name,
	)
	if err != nil {
		return err
	}
	s.status.Created = true
	s.ID = &containerCreateResponse.ID
	return nil
}

// Start starts the container. It returns an error if something goes wrong
func (c *Container) Start() error {
	// Start the container
	err := cli.ContainerStart(
		c.ctx,
		*c.ID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return err
	}
	c.status.Started = true
	return nil
}

// Stop stops the container and returns an error if something goes wrong
func (c *Container) Stop() error {
	err := cli.ContainerStop(c.ctx, *c.ID, container.StopOptions{Timeout: c.config.StopTimeout})
	if err != nil {
		return err
	}
	c.status.Started = false
	return nil
}

// Remove removes the container and returns an error if something goes wrong
func (c Container) Remove() error {
	err := cli.ContainerRemove(
		c.ctx,
		*c.ID,
		types.ContainerRemoveOptions{},
	)
	if err != nil {
		return err
	}
	c.status.Created = false
	c.status.Removed = true
	return nil
}

// Copyto copy the given slice of files to the target directory.
// It returns an error if something goes wrong
func (c Container) CopyTo(files []*os.File, targetDir *string) error {
	tarBuffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(tarBuffer)
	for _, file := range files {
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		header := &tar.Header{
			Name: fileInfo.Name(),
			Size: fileInfo.Size(),
			Mode: int64(fileInfo.Mode()),
		}
		err = tarWriter.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return err
		}
	}
	err := tarWriter.Close()
	if err != nil {
		return err
	}
	tdir := "/"
	if targetDir != nil {
		tdir = *targetDir
	}
	// Copy the tar archive to the container
	err = cli.CopyToContainer(
		c.ctx,
		*c.ID,
		tdir,
		tarBuffer,
		types.CopyToContainerOptions{CopyUIDGID: true},
	)
	return err
}
