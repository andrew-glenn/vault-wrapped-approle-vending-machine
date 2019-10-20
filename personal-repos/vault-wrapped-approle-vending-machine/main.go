package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	dockerApi "github.com/fsouza/go-dockerclient"
)

var dockerHost = flag.String("docker-host", getopt("DOCKER_HOST", "unix:///var/run/docker.sock"), "Address for the Docker daemon")

func getopt(name, def string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	return def
}
func readEventStream(events chan *dockerApi.APIEvents) {
	for msg := range events {
		fmt.Println("Event received!")
		fmt.Printf("%+v\n", msg)
	}
}

func watch_events(docker *dockerApi.Client) {

	events := make(chan *dockerApi.APIEvents)
	docker.AddEventListener(events)

	go readEventStream(events)
}

func main() {

	flag.Parse()

	log.Println("Listening for container events...")

	docker, err := dockerApi.NewClient(*dockerHost)

	if err != nil {
		log.Fatal(err)
	}

	watch_events(docker)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

forever:
	for {
		select {
		case <-sig:
			log.Println("signal received, stopping")
			break forever
		}
	}
}
