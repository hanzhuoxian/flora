package clientcmd

import "time"

type Server struct {
	LocationOfOrigin string
	Timeout          time.Duration
}

type ClientConfig struct{}
