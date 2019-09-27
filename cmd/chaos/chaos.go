package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/golang/glog"
	"github.com/livepeer/swarm-chaos/internal/engine"
	"github.com/livepeer/swarm-chaos/internal/engine/drivers/docker"
	"github.com/livepeer/swarm-chaos/internal/model"
)

func main() {
	flag.Set("logtostderr", "true")
	intMin := flag.String("int_min", "", "Interval, min")
	intMax := flag.String("int_max", "", "Interval, max")
	fKey := flag.String("f_key", "", "Label key")
	fVal := flag.String("f_val", "", "Label val")
	agent := flag.String("agent", "tcp://localhost:9001", "URL of the agent")
	server := flag.Bool("server", false, "Start in server mode")
	version := flag.Bool("version", false, "Print out the version")

	flag.Parse()

	if *version {
		fmt.Println("Swarm Chaos Version: " + model.SwarmChaosVersion)
		fmt.Printf("Golang runtime version: %s %s\n", runtime.Compiler, runtime.Version())
		fmt.Printf("Architecture: %s\n", runtime.GOARCH)
		fmt.Printf("Operating system: %s\n", runtime.GOOS)
		return
	}

	if *agent != "" {
		docker.AgentHost = *agent
	}

	if *server {
		dp, err := docker.NewDockerPlayground()
		if err != nil {
			panic(err)
		}
		scheduler := engine.NewScheduler(dp)
		server := engine.NewServer(scheduler)
		server.StartServer()
		return
	}
	if *intMin == "" {
		glog.Info("int_min must be specified")
		return
	}
	if *intMax == "" {
		glog.Info("int_max must be specified")
		return
	}
	if *fKey == "" {
		glog.Info("f_key must be specified")
		return
	}
	if *fVal == "" {
		glog.Info("f_val must be specified")
		return
	}
	dp, err := docker.NewDockerPlayground()
	if err != nil {
		panic(err)
	}
	scheduler := engine.NewScheduler(dp)
	err = scheduler.ScheduleTask(*intMin, *intMax, model.OperationTypeDestroy, *fKey, *fVal)
	if err != nil {
		panic(err)
	}
	scheduler.StartTasks()

	runtime.Goexit()

	/*
		engine := engine.NewChaosEngine()
		dp, err := docker.NewDockerPlayground()
		if err != nil {
			panic(err)
		}
		engine.AddPlayground(dp)
		ents, err := engine.EntitiesByLabel("type", "transcoder")
		// ents, err := engine.Entities()
		if err != nil {
			panic(err)
		}
		for _, ent := range ents {
			fmt.Printf("Container named %s\n", ent.Name())
			fmt.Printf("Container %s has labels:\n", ent.Name())
			fmt.Printf("%+v\n", ent.Labels())
		}
		if len(ents) > 0 {
			ent := ents[0]
			fmt.Printf("Killing %s\n", ent.Name())
			// err := ent.Do(model.OperationTypeDestroy)
			// if err != nil {
			// 	panic(err)
			// }
		}
	*/

	/*
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv)
		// cli, err := client.NewClientWithOpts(client.W)
		if err != nil {
			panic(err)
		}
		cli.NegotiateAPIVersion(ctx)
		swa, err := cli.SwarmInspect(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("swarm id: %s\n", swa.ID)
		fmt.Printf("Swarm info: %+v\n", swa)

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		// Вывод всех идентификаторов контейнеров
		for _, container := range containers {
			fmt.Println(container.ID)
			fmt.Printf("state: %s status: %s\n", container.State, container.Status)
			for _, node := range container.Names {
				fmt.Println(node)
			}
		}
	*/

	/*
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}
		cli.NegotiateAPIVersion(ctx)

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		// Вывод всех идентификаторов контейнеров
		for _, container := range containers {
			fmt.Println(container.ID)
			fmt.Printf("state: %s status: %s\n", container.State, container.Status)
			for _, node := range container.Names {
				fmt.Println(node)
			}
		}

			reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
			if err != nil {
				panic(err)
			}
			io.Copy(os.Stdout, reader)

			resp, err := cli.ContainerCreate(ctx, &container.Config{
				Image: "alpine",
				Cmd:   []string{"echo", "hello world"},
			}, nil, nil, "")
			if err != nil {
				panic(err)
			}

			if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				panic(err)
			}

			statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
			select {
			case err := <-errCh:
				if err != nil {
					panic(err)
				}
			case <-statusCh:
			}

			out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
			if err != nil {
				panic(err)
			}

			stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	*/
}
