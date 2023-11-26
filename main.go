package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hekmon/transmissionrpc/v2"
)

type GluetunResponse struct {
	Port uint16
}

var (
	transmissionHostname = flag.String("transmission-hostname", "127.0.0.1", "transmission hostname")
	transmissionPort     = flag.Int("transmission-port", 9091, "transmission port")
	transmissionUsername = os.Getenv("TRANSMISSION_USERNAME")
	transmissionPassword = os.Getenv("TRANSMISSION_PASSWORD")

	gluetunHostname = getEnv("GLUETUN_HOSTNAME", "127.0.0.1")
	gluetunPort     = getEnv("GLUETUN_PORT", "8000")

	initialDelayStr  = getEnv("INITIAL_DELAY", "5s")
	checkIntervalStr = getEnv("CHECK_INTERVAL", "1m")
	errorIntervalStr = getEnv("ERROR_INTERVAL", "5s")
)

func init() {
	flag.Parse()
}

func main() {
	initialDelay, _ := time.ParseDuration(initialDelayStr)
	checkInterval, _ := time.ParseDuration(checkIntervalStr)
	errorInterval, _ := time.ParseDuration(errorIntervalStr)
	previousExternalPort := uint16(0)
	gluetunPortApi := fmt.Sprintf("http://%s:%s/v1/openvpn/portforwarded", gluetunHostname, gluetunPort)
	errorCount := 0
	maxErrorCount := 5

	time.Sleep(initialDelay)

	httpClient := resty.New()

	transmissionClient, err := transmissionrpc.New(*transmissionHostname, transmissionUsername, transmissionPassword, &transmissionrpc.AdvancedConfig{
		Port: uint16(*transmissionPort),
	})
	if err != nil {
		log.Fatalf("failed to create transmission client: %v", err)
	}

	for {
		portMapping := &GluetunResponse{}
		resp, err := httpClient.R().
			SetResult(portMapping).
			ForceContentType("application/json").
			Get(gluetunPortApi)

		if err != nil || resp.IsError() {
			log.Fatalf("failed to fetch port mapping from gluetun: %v, %d", err, resp.StatusCode())
		}

		if portMapping.Port == 0 {
			err = errors.New("empty port")
			log.Printf("new port is not yet assigned, %v", err)
		} else if portMapping.Port == previousExternalPort {
			log.Printf("external port is unchanged: %d", portMapping.Port)
		} else {
			log.Printf("external port changed to: %d", portMapping.Port)

			transmissionPeerPort := int64(portMapping.Port)
			err = transmissionClient.SessionArgumentsSet(context.Background(), transmissionrpc.SessionArguments{
				PeerPort: &transmissionPeerPort,
			})
			if err != nil {
				log.Fatalf("failed to set transmission peer port: %v", err)
			} else {
				previousExternalPort = portMapping.Port
				log.Printf("updated transmission peer port to: %d", transmissionPeerPort)
			}
		}

		if err != nil && errorCount < maxErrorCount {
			time.Sleep(errorInterval)
			errorCount++
		} else {
			time.Sleep(checkInterval)
			errorCount = 0
		}
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
