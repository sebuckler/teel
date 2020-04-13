package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewRoot(v string) *cobra.Command {
	root := &cobra.Command{
		Use:     "teel",
		Short:   "Teel is a static blog management system",
		Version: v,
	}

	root.SetVersionTemplate(fmt.Sprintf("Teel Static Blog Management System %s\n", v))

	return root
}
