// +build linux darwin freebsd netbsd openbsd
// +build !nofuse

package blns

import (
	core "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	mount "github.com/ipfs/go-ipfs/fuse/mount"
)

// Mount mounts blns at a given location, and returns a mount.Mount instance.
func Mount(link *core.IpfsNode, blnsmp, linkmp string) (mount.Mount, error) {
	coreApi, err := coreapi.NewCoreAPI(link)
	if err != nil {
		return nil, err
	}

	cfg, err := link.Repo.Config()
	if err != nil {
		return nil, err
	}

	allow_other := cfg.Mounts.FuseAllowOther

	fsys, err := NewFileSystem(link.Context(), coreApi, linkmp, blnsmp)
	if err != nil {
		return nil, err
	}

	return mount.NewMount(link.Process, fsys, blnsmp, allow_other)
}
