package clientcmd

import "time"

type Server struct {
	LocationOfOrigin string
	Timeout          time.Duration
}
