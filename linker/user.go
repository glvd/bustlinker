package linker

import (
	"context"
	"errors"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/libp2p/go-libp2p-core/peer"
)

type user struct {
	node      *core.IpfsNode
	peerCache PeerCache
}

func (u *user) Subscribe(ctx context.Context, id string) error {
	api, err := coreapi.NewCoreAPI(u.node)
	if err != nil {
		return err
	}
	pid, err := peer.IDFromString(id)

	if err != nil {
		return err
	}
	address, b := u.peerCache.GetAddress(pid)
	if !b {
		return errors.New("remote address not found")
	}
	if err := api.Swarm().Connect(ctx, address); err != nil {
		return err
	}
	return nil
}

func (u *user) Describe(id string) {

}
