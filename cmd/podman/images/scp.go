package images

import (
	"os"
	"strings"

	"github.com/containers/podman/v4/cmd/podman/common"
	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/spf13/cobra"
)

var (
	saveScpDescription = `Securely copy an image from one host to another.`
	imageScpCommand    = &cobra.Command{
		Use: "scp [options] IMAGE [HOST::]",
		Annotations: map[string]string{
			registry.UnshareNSRequired: "",
			registry.ParentNSRequired:  "",
		},
		Long:              saveScpDescription,
		Short:             "securely copy images",
		RunE:              scp,
		Args:              cobra.RangeArgs(1, 2),
		ValidArgsFunction: common.AutocompleteScp,
		Example:           `podman image scp myimage:latest otherhost::`,
	}
)

var (
	parentFlags []string
	quiet       bool
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: imageScpCommand,
		Parent:  imageCmd,
	})
	scpFlags(imageScpCommand)
}

func scpFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.BoolVarP(&quiet, "quiet", "q", false, "Suppress the output")
}

func scp(cmd *cobra.Command, args []string) (finalErr error) {
	var (
		err error
	)
	for i, val := range os.Args {
		if val == "image" {
			break
		}
		if i == 0 {
			continue
		}
		if strings.Contains(val, "CIRRUS") { // need to skip CIRRUS flags for testing suite purposes
			continue
		}
		parentFlags = append(parentFlags, val)
	}

	src := args[0]
	dst := ""
	if len(args) > 1 {
		dst = args[1]
	}

	err = registry.ImageEngine().Scp(registry.Context(), src, dst, parentFlags, quiet)
	if err != nil {
		return err
	}

	return nil
}
