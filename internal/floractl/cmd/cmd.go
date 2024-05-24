package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	cliflag "github.com/hanzhuoxian/flora/pkg/cli/flag"
	"github.com/hanzhuoxian/flora/pkg/log"
)

func NewDefaultCommand() *cobra.Command {
	return NewCommand(os.Stdin, os.Stdout, os.Stderr)
}

func NewCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	log.Init(log.NewOptions())

	cmds := &cobra.Command{
		Use:   "floractl",
		Short: "floractl is a command line tool for managing the Flora project",
		Long:  "floractl is a command line tool for managing the Flora project",
		Run:   runHelp,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initProfiling()
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return flushProfiling()
		},
	}

	flags := cmds.PersistentFlags()

	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc)

	addProfilingFlags(flags)

	// configFlags := clioptions.NewConfigFlags(true)

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
