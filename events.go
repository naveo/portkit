package main

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func DockerController() (c *DockerCli, err error) {
	c = new(DockerCli)
	c.cli, err = client.NewClientWithOpts(client.WithHost(config.NaveoHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *DockerCli) EventFilter() (<-chan events.Message, <-chan error) {
	// Docker events to listen for events
	filter := filters.NewArgs()
	filter.Add("type", "container")
	filter.Add("event", "start")
	filter.Add("event", "die")

	eventChannel, errorChannel := c.cli.Events(context.Background(), types.EventsOptions{
		Filters: filter,
	})

	return eventChannel, errorChannel
}

func (c *DockerCli) ContainerInspector(containerEvent string) (types.ContainerJSON, error) {
	state, err := c.cli.ContainerInspect(context.Background(), containerEvent)
	return state, err
}

func DockerEventListener() (*DockerCli, <-chan events.Message, <-chan error) {
	cli, err := DockerController()
	if err != nil {
		panic(err)
	}

	eventChannel, errorChannel := cli.EventFilter()

	return cli, eventChannel, errorChannel
}

func DockerEvents() {
	cli, eventChannel, errorChannel := DockerEventListener()

	for {
		select {
		case err := <-errorChannel:
			log.Fatalf("error reading events channel %v:", err)
		case containerEvent := <-eventChannel:
			state, err := cli.ContainerInspector(containerEvent.ID)
			if err != nil {
				log.Fatalf("failed getting state for container ID: %s error: %v", containerEvent.ID, err)
			}
			go PortForwardParser(state, containerEvent.ID)
		}
	}
}

// PortForwardParser iterates over all requested port map index like 80/tcp and over
// all ports listed within the index ignoring IPv6 format to call the PortForwardMapper
// on all requested ports.
func PortForwardParser(state types.ContainerJSON, id string) {
	for index := range state.NetworkSettings.NetworkSettingsBase.Ports {
		for ports := range state.NetworkSettings.NetworkSettingsBase.Ports[index] {
			if !strings.Contains(state.NetworkSettings.NetworkSettingsBase.Ports[index][ports].HostIP, ":") {
				port := state.NetworkSettings.NetworkSettingsBase.Ports[index][ports].HostPort
				log.Printf("initiate port forward for cid %v port %v", id, port)
				go PortMapper(id, port)
			}
		}
	}
}
