package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hekmon/transmissionrpc/v3"
)

type GluetunResponse struct {
	Port uint16
}

var (
	// Transmission
	transmissionHostname = flag.String("transmission-hostname", "127.0.0.1", "transmission hostname")
	transmissionPort     = flag.Int("transmission-port", 9091, "transmission port")
	transmissionUsername = os.Getenv("TRANSMISSION_USERNAME")
	transmissionPassword = os.Getenv("TRANSMISSION_PASSWORD")
	transmissionProtocol = getEnv("TRANSMISSION_PROTOCOL", "http")

	// Gluetun
	gluetunProtocol     = getEnv("GLUETUN_PROTOCOL", "http")
	gluetunHostname     = getEnv("GLUETUN_HOSTNAME", "127.0.0.1")
	gluetunPort         = getEnv("GLUETUN_PORT", "8000")
	gluetunEndpoint     = getEnv("GLUETUN_ENDPOINT", "/v1/portforward")
	gluetunAuthType     = getEnv("GLUETUN_AUTH_TYPE", "none")
	gluetunAuthUsername = os.Getenv("GLUETUN_AUTH_USERNAME")
	gluetunAuthPassword = os.Getenv("GLUETUN_AUTH_PASSWORD")
	gluetunAuthAPIKey   = os.Getenv("GLUETUN_AUTH_API_KEY")

	// Control flow
	initialDelayStr  = getEnv("INITIAL_DELAY", "5s")
	checkIntervalStr = getEnv("CHECK_INTERVAL", "1m")
	errorIntervalStr = getEnv("ERROR_INTERVAL", "5s")

	gluetunAuthTypeBasic        = "basic"
	gluetunAuthTypeAPIKey       = "apikey"
	gluetunAuthTypeAPIKeyHeader = "X-API-Key"
)

func init() {
	flag.Parse()
}

func main() {
	initialDelay, _ := time.ParseDuration(initialDelayStr)
	checkInterval, _ := time.ParseDuration(checkIntervalStr)
	errorInterval, _ := time.ParseDuration(errorIntervalStr)
	previousExternalPort := uint16(0)

	if !strings.HasPrefix(gluetunEndpoint, "/") {
		gluetunEndpoint = fmt.Sprintf("/%s", gluetunEndpoint)
	}

	gluetunPortApi := fmt.Sprintf("%s://%s:%s%s", gluetunProtocol, gluetunHostname, gluetunPort, gluetunEndpoint)
	errorCount := 0
	maxErrorCount := 5

	time.Sleep(initialDelay)

	httpClient := resty.New()
	authInfo := ""
	if transmissionUsername != "" && transmissionPassword != "" {
		authInfo = fmt.Sprintf("%s@", url.UserPassword(transmissionUsername, transmissionPassword).String())
	}

	endpoint, err := url.Parse(fmt.Sprintf("%s://%s%s:%d/transmission/rpc", transmissionProtocol, authInfo, *transmissionHostname, *transmissionPort))
	if err != nil {
		log.Fatalf("failed to parse transmission endpoint: %v", err)
	}

	transmissionClient, err := transmissionrpc.New(endpoint, nil)
	if err != nil {
		log.Fatalf("failed to create transmission client: %v", err)
	}

	log.Printf("fetching port mapping from gluetun using auth type: %s", gluetunAuthType)

	for {
		portMapping := &GluetunResponse{}
		req := httpClient.R().
			SetResult(portMapping).
			ForceContentType("application/json")

		switch gluetunAuthType {
		case gluetunAuthTypeBasic:
			req.SetBasicAuth(gluetunAuthUsername, gluetunAuthPassword)
		case gluetunAuthTypeAPIKey:
			req.SetHeader(gluetunAuthTypeAPIKeyHeader, gluetunAuthAPIKey)
		}

		resp, err := req.Get(gluetunPortApi)
		if err != nil || resp.IsError() {
			log.Printf("failed to fetch port mapping from gluetun: %v, %d", err, resp.StatusCode())
		}

		if portMapping.Port == 0 {
			err = errors.New("empty port")
			log.Printf("new port is not yet assigned, %v", err)
		} else if portMapping.Port != previousExternalPort {
			log.Printf("external port changed to: %d", portMapping.Port)

			transmissionPeerPort := int64(portMapping.Port)
			err = transmissionClient.SessionArgumentsSet(context.Background(), transmissionrpc.SessionArguments{
				PeerPort: &transmissionPeerPort,
			})

			if err != nil {
				log.Printf("failed to set transmission peer port: %v", err)
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

/**
 * Return fallback fallback if environment variable is not set or empty
 */
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
