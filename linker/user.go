package linker

import (
	"context"
	"github.com/ipfs/go-ipfs/core"
)

type user struct {
	node *core.IpfsNode
}

func (u *user) Subscribe(ctx context.Context, id string) error {
	return nil
}

func (u *user) Describe(id string) {

}
