package resolve

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ipfs/go-path"

	"github.com/ipfs/go-ipfs/namesys"
)

// ErrNoNamesys is an explicit error for when an LINK node doesn't
// (yet) have a name system
var ErrNoNamesys = errors.New(
	"core/resolve: no Namesys on IpfsNode - can't resolve blns entry")

// ResolveBLNS resolves /blns paths
func ResolveBLNS(ctx context.Context, nsys namesys.NameSystem, p path.Path) (path.Path, error) {
	if strings.HasPrefix(p.String(), "/blns/") {
		// TODO(cryptix): we should be able to query the local cache for the path
		if nsys == nil {
			return "", ErrNoNamesys
		}

		seg := p.Segments()

		if len(seg) < 2 || seg[1] == "" { // just "/<protocol/>" without further segments
			err := fmt.Errorf("invalid path %q: blns path missing BLNS ID", p)
			return "", err
		}

		extensions := seg[2:]
		resolvable, err := path.FromSegments("/", seg[0], seg[1])
		if err != nil {
			return "", err
		}

		respath, err := nsys.Resolve(ctx, resolvable.String())
		if err != nil {
			return "", err
		}

		segments := append(respath.Segments(), extensions...)
		p, err = path.FromSegments("/", segments...)
		if err != nil {
			return "", err
		}
	}
	return p, nil
}
