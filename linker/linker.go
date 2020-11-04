package linker

import (
	"context"
	"encoding/base64"
	"fmt"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/linker/config"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"sync"
)

const Version = "0.0.1"
const LinkPeers = "/link" + "/peers/" + Version
const LinkAddress = "/link" + "/address/" + Version
const LinkHash = "/link" + "/hash/" + Version

var protocols = []string{
	LinkPeers,
	LinkAddress,
}

var NewLine = []byte{'\n'}

type Linker interface {
	Start(node *core.IpfsNode) error
	//plugin.Plugin
	//plugin.PluginDaemonInternal
}

type link struct {
	ctx         context.Context
	cfg         *config.Config
	node        *core.IpfsNode
	failedCount map[peer.ID]int64
	failedLock  *sync.RWMutex
	pinning     Pinning
	repo        string
}

func (l *link) newLinkPeersHandle() (protocol.ID, func(stream network.Stream)) {
	return LinkPeers, func(stream network.Stream) {
		log.Debug("link peer called")
		var err error
		defer stream.Close()
		remoteID := stream.Conn().RemotePeer()

		peers := l.node.PeerHost.Network().Peers()
		log.Infow("get all peers", "total", len(peers))
		for _, peer := range peers {
			info := l.node.Peerstore.PeerInfo(peer)
			json, _ := info.MarshalJSON()
			_, err = stream.Write(json)
			if err != nil {
				log.Debugw("stream write error", "error", err)
				return
			}
			_, err = stream.Write(NewLine)
			if err != nil {
				log.Debugw("stream write error", "error", err)
				return
			}
			log.Debugw("peer sent success", "to", remoteID, stream.Conn().RemoteMultiaddr().String(), "addr", info.String())
		}
	}
}

func (l *link) newLinkHashHandle() (protocol.ID, func(stream network.Stream)) {
	return LinkHash, func(stream network.Stream) {
		log.Debug("link hash called")
		var err error
		defer stream.Close()
		for _, peer := range l.pinning.Get() {
			_, err = stream.Write([]byte(peer))
			if err != nil {
				log.Debugw("stream write error", "error", err)
				return
			}
			_, err = stream.Write(NewLine)
			if err != nil {
				log.Debugw("stream write error", "error", err)
				return
			}
		}
	}
}

func (l *link) registerHandle() {
	l.node.PeerHost.SetStreamHandler(l.newLinkPeersHandle())
	l.node.PeerHost.SetStreamHandler(l.newLinkHashHandle())
}

func (l *link) Start(node *core.IpfsNode) error {
	fmt.Println("Link start")
	l.node = node

	l.pinning = newPinning(l.node)

	l.registerHandle()
	return nil
}

func New(repo string, cfg interface{}) (Linker, error) {
	v, b := cfg.(*config.Config)
	if cfg == nil || !b {
		v = config.InitConfig(repo)
		err := config.StoreConfig(repo, v)
		if err != nil {
			return nil, err
		}
	}
	return &link{
		ctx:         context.TODO(),
		repo:        repo,
		cfg:         v,
		failedCount: make(map[peer.ID]int64),
		failedLock:  &sync.RWMutex{},
	}, nil
}

func LinkKey(id peer.ID) ds.Key {
	return ds.NewKey("/link/" + base64.RawStdEncoding.EncodeToString([]byte(id)))
}

var _ Linker = &link{}
