package name

import (
	"fmt"
	"io"
	"strings"

	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
	ke "github.com/ipfs/go-ipfs/core/commands/keyencode"
	"github.com/libp2p/go-libp2p-core/peer"
	record "github.com/libp2p/go-libp2p-record"
)

type blnsPubsubState struct {
	Enabled bool
}

type blnsPubsubCancel struct {
	Canceled bool
}

type stringList struct {
	Strings []string
}

// IpnsPubsubCmd is the subcommand that allows us to manage the BLNS pubsub system
var IpnsPubsubCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "BLNS pubsub management",
		ShortDescription: `
Manage and inspect the state of the BLNS pubsub resolver.

Note: this command is experimental and subject to change as the system is refined
`,
	},
	Subcommands: map[string]*cmds.Command{
		"state":  blnspsStateCmd,
		"subs":   blnspsSubsCmd,
		"cancel": blnspsCancelCmd,
	},
}

var blnspsStateCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Query the state of BLNS pubsub",
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &blnsPubsubState{n.PSRouter != nil})
	},
	Type: blnsPubsubState{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, ips *blnsPubsubState) error {
			var state string
			if ips.Enabled {
				state = "enabled"
			} else {
				state = "disabled"
			}

			_, err := fmt.Fprintln(w, state)
			return err
		}),
	},
}

var blnspsSubsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Show current name subscriptions",
	},
	Options: []cmds.Option{
		ke.OptionBLNSBase,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		keyEnc, err := ke.KeyEncoderFromString(req.Options[ke.OptionBLNSBase.Name()].(string))
		if err != nil {
			return err
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		if n.PSRouter == nil {
			return cmds.Errorf(cmds.ErrClient, "BLNS pubsub subsystem is not enabled")
		}
		var paths []string
		for _, key := range n.PSRouter.GetSubscriptions() {
			ns, k, err := record.SplitKey(key)
			if err != nil || ns != "blns" {
				// Not necessarily an error.
				continue
			}
			pid, err := peer.IDFromString(k)
			if err != nil {
				log.Errorf("blns key not a valid peer ID: %s", err)
				continue
			}
			paths = append(paths, "/blns/"+keyEnc.FormatID(pid))
		}

		return cmds.EmitOnce(res, &stringList{paths})
	},
	Type: stringList{},
	Encoders: cmds.EncoderMap{
		cmds.Text: stringListEncoder(),
	},
}

var blnspsCancelCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Cancel a name subscription",
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		if n.PSRouter == nil {
			return cmds.Errorf(cmds.ErrClient, "BLNS pubsub subsystem is not enabled")
		}

		name := req.Arguments[0]
		name = strings.TrimPrefix(name, "/blns/")
		pid, err := peer.Decode(name)
		if err != nil {
			return cmds.Errorf(cmds.ErrClient, err.Error())
		}

		ok, err := n.PSRouter.Cancel("/blns/" + string(pid))
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, &blnsPubsubCancel{ok})
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("name", true, false, "Name to cancel the subscription for."),
	},
	Type: blnsPubsubCancel{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, ipc *blnsPubsubCancel) error {
			var state string
			if ipc.Canceled {
				state = "canceled"
			} else {
				state = "no subscription"
			}

			_, err := fmt.Fprintln(w, state)
			return err
		}),
	},
}

func stringListEncoder() cmds.EncoderFunc {
	return cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, list *stringList) error {
		for _, s := range list.Strings {
			_, err := fmt.Fprintln(w, s)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
