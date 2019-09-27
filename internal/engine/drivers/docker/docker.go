package docker

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
	"github.com/livepeer/swarm-chaos/internal/model"
)

// AgentHost is url of the agent
var AgentHost = "tcp://localhost:9001"

type (
	DockerPlayground struct {
		client *client.Client
		nodes  []swarm.Node
	}

	dockerContainer struct {
		container types.Container
		dp        *DockerPlayground
		client    *client.Client
	}
)

var operationTimeout = 2 * time.Second

func newPortainerAgentClient(host string, nodeName string) (*client.Client, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	httpCli := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	headers := map[string]string{
		"X-PortainerAgent-PublicKey": "3059301306072a8648ce3d020106082a8648ce3d0301070342000442afbfddf2c43e5c9175762c57015275892da81d518eb5700810ae4278600ef22779d2701a2c1039edc4cfc532840178d861d4ea1550c63f2ee920a1aec18dc9",
		"X-PortainerAgent-Signature": "lz4dAYZ4Xqw4Q+1yVQIVUKyDCQdg8hcn7ijIMTSVl1dbOTxEQtbBdPGkn4hMcW6R+vhiyTCft+ANYhfDHBKzyQ",
	}

	if nodeName != "" {
		headers["X-PortainerAgent-Target"] = nodeName
	}
	fmt.Printf("Using headers: %+v\n", headers)

	return client.NewClientWithOpts(
		client.WithHost(host),
		client.WithHTTPClient(httpCli),
		client.WithHTTPHeaders(headers),
	)
}

// NewDockerPlayground creates a new docker driver
func NewDockerPlayground() (*DockerPlayground, error) {
	ctx := context.Background()
	// cli, err := client.NewClientWithOpts(client.FromEnv)
	cli, err := newPortainerAgentClient(AgentHost, "")
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)
	dp := &DockerPlayground{
		client: cli,
	}
	dp.refreshNodes()

	return dp, nil
}

// Entities returns a list of all the entities
func (dp *DockerPlayground) Entities() ([]model.Entity, error) {
	res := make([]model.Entity, 0)
	containers, err := dp.client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		// fmt.Printf("Got container info ---------------------------:\n")
		// fmt.Printf("%+v\n", container)
		dc := &dockerContainer{
			container: container,
			dp:        dp,
			// client:    dp.client,
		}
		res = append(res, dc)
	}

	return res, nil
}

func (dp *DockerPlayground) refreshNodes() {
	ctx := context.Background()
	nodes, err := dp.client.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		glog.Infof("Error getting nodes list: %v", err)
		return
	}
	dp.nodes = nodes
	for _, node := range dp.nodes {
		glog.Infof("Node %s has id %s\n", node.Description.Hostname, node.ID)
	}
}

func (dp *DockerPlayground) getNodeNameFromLabels(labels map[string]string) string {
	// glog.Infof("get node name form labels: %+v\n", labels)
	// glog.Infof("nodes list: %+v\n", dp.nodes)
	// glog.Infof("nodes number %d\n", len(dp.nodes))
	if len(dp.nodes) == 0 {
		dp.refreshNodes()
	}
	if id, has := labels["com.docker.swarm.node.id"]; has {
		glog.Infof("node id: %s, has: %v\n", id, has)
		for _, node := range dp.nodes {
			if node.ID == id {
				return node.Description.Hostname
			}
		}
	}
	return ""
}

func (dc *dockerContainer) Name() string {
	name := dc.container.ID
	if len(dc.container.Names) > 0 {
		name = dc.container.Names[0]
		// name = name + ":" + strings.Join(dc.container.Names, "|")
	}

	return name
}

func (dc *dockerContainer) Labels() map[string]string {
	return dc.container.Labels
}

func (dc *dockerContainer) Childs() []model.Entity {
	return nil
}

func (dc *dockerContainer) Type() model.EntityType {
	return model.EntityTypeContainer
}

func (dc *dockerContainer) Do(operation model.OperationType) error {
	client := dc.getClient()
	switch operation {
	case model.OperationTypeDestroy:
		err := client.ContainerRemove(context.Background(), dc.container.ID, types.ContainerRemoveOptions{Force: true})
		return err
	case model.OperationTypePause:
		err := client.ContainerPause(context.Background(), dc.container.ID)
		return err
	case model.OperationTypeResume:
		err := client.ContainerUnpause(context.Background(), dc.container.ID)
		return err
	case model.OperationTypeStop:
		err := client.ContainerStop(context.Background(), dc.container.ID, &operationTimeout)
		return err
	case model.OperationTypeStart:
		err := client.ContainerStart(context.Background(), dc.container.ID, types.ContainerStartOptions{})
		return err
	}
	return nil
}

func (dc *dockerContainer) Status() (model.StatusType, error) {
	// dc.container.Status
	client := dc.getClient()
	j, err := client.ContainerInspect(context.Background(), dc.container.ID)
	if err != nil {
		return 0, err
	}
	// Status     string // String representation of the container state. Can be one of "created", "running", "paused", "restarting", "removing", "exited", or "dead"
	switch j.State.Status {
	case "created", "running", "restarting":
		return model.StatusTypeWorking, nil
	case "paused":
		return model.StatusTypePaused, nil
	case "exited", "dead", "removing":
		return model.StatusTypeDestroyed, nil
	}
	return model.StatusTypeWorking, nil
}

func (dc *dockerContainer) getClient() *client.Client {
	if dc.client != nil {
		return dc.client
	}
	nodeName := dc.dp.getNodeNameFromLabels(dc.container.Labels)
	glog.Infof("Got node name from labels: %s\n", nodeName)
	if nodeName == "" {
		dc.client = dc.dp.client
	} else {
		cli, err := newPortainerAgentClient(AgentHost, nodeName)
		if err != nil {
			panic(err)
		}
		cli.NegotiateAPIVersion(context.Background())
		dc.client = cli
	}
	return dc.client
}
