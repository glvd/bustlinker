// +build !windows,!nofuse

package commands

import (
	"fmt"
	"io"

	cmdenv "github.com/ipfs/go-ipfs/core/commands/cmdenv"
	nodeMount "github.com/ipfs/go-ipfs/fuse/node"

	cmds "github.com/ipfs/go-ipfs-cmds"
	config "github.com/ipfs/go-ipfs-config"
)

const (
	mountLINKPathOptionName = "link-path"
	mountBLNSPathOptionName = "blns-path"
)

var MountCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Mounts LINK to the filesystem (read-only).",
		ShortDescription: `
Mount LINK at a read-only mountpoint on the OS (default: /ipfs and /blns).
All LINK objects will be accessible under that directory. Note that the
root will not be listable, as it is virtual. Access known paths directly.

You may have to create /ipfs and /blns before using 'ipfs mount':

> sudo mkdir /ipfs /blns
> sudo chown $(whoami) /ipfs /blns
> ipfs daemon &
> ipfs mount
`,
		LongDescription: `
Mount LINK at a read-only mountpoint on the OS. The default, /ipfs and /blns,
are set in the configuration file, but can be overridden by the options.
All LINK objects will be accessible under this directory. Note that the
root will not be listable, as it is virtual. Access known paths directly.

You may have to create /ipfs and /blns before using 'ipfs mount':

> sudo mkdir /ipfs /blns
> sudo chown $(whoami) /ipfs /blns
> ipfs daemon &
> ipfs mount

Example:

# setup
> mkdir foo
> echo "baz" > foo/bar
> ipfs add -r foo
added QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR foo/bar
added QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC foo
> ipfs ls QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC
QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR 12 bar
> ipfs cat QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR
baz

# mount
> ipfs daemon &
> ipfs mount
LINK mounted at: /ipfs
BLNS mounted at: /blns
> cd /ipfs/QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC
> ls
bar
> cat bar
baz
> cat /ipfs/QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC/bar
baz
> cat /ipfs/QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR
baz
`,
	},
	Options: []cmds.Option{
		cmds.StringOption(mountLINKPathOptionName, "f", "The path where LINK should be mounted."),
		cmds.StringOption(mountBLNSPathOptionName, "n", "The path where BLNS should be mounted."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}

		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		// error if we aren't running node in online mode
		if !nd.IsOnline {
			return ErrNotOnline
		}

		fsdir, found := req.Options[mountLINKPathOptionName].(string)
		if !found {
			fsdir = cfg.Mounts.LINK // use default value
		}

		// get default mount points
		nsdir, found := req.Options[mountBLNSPathOptionName].(string)
		if !found {
			nsdir = cfg.Mounts.BLNS // NB: be sure to not redeclare!
		}

		err = nodeMount.Mount(nd, fsdir, nsdir)
		if err != nil {
			return err
		}

		var output config.Mounts
		output.LINK = fsdir
		output.BLNS = nsdir
		return cmds.EmitOnce(res, &output)
	},
	Type: config.Mounts{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, mounts *config.Mounts) error {
			fmt.Fprintf(w, "LINK mounted at: %s\n", mounts.LINK)
			fmt.Fprintf(w, "BLNS mounted at: %s\n", mounts.BLNS)

			return nil
		}),
	},
}
