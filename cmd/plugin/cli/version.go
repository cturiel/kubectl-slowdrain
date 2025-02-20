package cli

import (
	"fmt"

	"github.com/cturiel/kubectl-slowdrain/pkg/version"
	"github.com/spf13/cobra"
)

// NewVersionCmd crea el comando `kubectl slowdrain version`
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the version of slowdrain",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("slowdrain version:", version.Version)
		},
	}
}
