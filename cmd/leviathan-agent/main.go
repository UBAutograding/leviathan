package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/UBAutograding/leviathan/internal/dockerclient"
	"github.com/UBAutograding/leviathan/internal/util"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)
	if log.GetLevel() == log.TraceLevel {
		log.SetReportCaller(true)
	}
	// TODO: For prod ensure logs are json formatted for ingest
	// log.SetFormatter(&log.JSONFormatter{})

	cli, err := client.NewEnvClient()
	// cli, err := dockerclient.NewSSHClient("yeager")
	if err != nil {
		log.Fatal("Failed to setup docker client")
	}

	// err = dockerclient.PullImage(cli, "ubautograding/autograding_image_2004")

	id, err := dockerclient.CreateNewContainer(cli, "ubautograding/autograding_image_2004")
	if err != nil {
		log.Fatal("Failed to create image")
	}

	err = dockerclient.CopyToContainer(cli, id, fmt.Sprintf("%s/tmp/sanitycheck/tmp/", util.UserHomeDir()))
	if err != nil {
		cleanup(cli, id)
	}

	err = dockerclient.StartContainer(cli, id)
	if err != nil {
		cleanup(cli, id)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go dockerclient.TailContainerLogs(ctx, cli, id)

	time.Sleep(10 * time.Second)
	cancel()

	cleanup(cli, id)
	dockerclient.ListContainer(cli)
}

func cleanup(c *client.Client, containerID string) {
	dockerclient.StopContainer(c, containerID)
	dockerclient.RemoveContainer(c, containerID, false, false)
	os.Exit(0)
}
