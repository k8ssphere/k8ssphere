package drone

import (
	"context"
	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
)

type Drone struct {
	Host  string
	Token string
}

/**
init
*/
func (d Drone) Init() drone.Client {
	config := new(oauth2.Config)
	auther := config.Client(
		context.Background(),
		&oauth2.Token{
			AccessToken: d.Token,
		},
	)
	return drone.NewClient(d.Host, auther)
}
