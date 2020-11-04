package linker

import (
	core "github.com/libp2p/go-libp2p-core"
	"sync"
)

type peerLink struct {
	lock  sync.RWMutex
	peers map[string]core.PeerAddrInfo
}
