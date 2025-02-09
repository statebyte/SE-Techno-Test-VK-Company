package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	ContainerID   string    `json:"container_id"`
	ContainerName string    `json:"container_name"`
	IPAddress     string    `json:"ip_address"`
	Status        string    `json:"status"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func main() {
	backendURL := getEnv("BACKEND_URL", "http://localhost:8080")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("Pinger запущен, будет ping каждые 10 секунд...")

	for {
		select {
		case <-ticker.C:
			if err := pingDockerAndSend(backendURL); err != nil {
				log.Printf("Ошибка pingDockerAndSend: %v", err)
			}
		}
	}
}

func pingDockerAndSend(backendURL string) error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err
	}

	var data []ContainerInfo
	for _, c := range containers {
		ip := ""
		if len(c.NetworkSettings.Networks) > 0 {
			for _, netConf := range c.NetworkSettings.Networks {
				ip = netConf.IPAddress
				break
			}
		}

		containerInfo := ContainerInfo{
			ContainerID:   c.ID[:12],
			ContainerName: c.Names[0],
			IPAddress:     ip,
			Status:        c.Status,
		}
		data = append(data, containerInfo)
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url := backendURL + "/api/containers"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("WARN: backend вернул статус %d", resp.StatusCode)
	}

	return nil
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
