package linker

import (
	"context"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/linker/data"
)

type user struct {
	node *core.IpfsNode
}

func (u *user) Subscribe(ctx context.Context, id string) error {
	return nil
}

func (u *user) Describe(id string) {

}

func (u *user) List() <-chan *data.User {
	userData := make(chan *data.User)

	return userData
}
