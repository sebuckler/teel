package commands

import (
	"fmt"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/spf13/cobra"
)

type CreateCommand struct {
	blogName       string
	blogScaffolder scaffolder.Scaffolder
	logWriter      logger.Logger
	targetDir      string
	*cobra.Command
}

func NewCreate(s scaffolder.Scaffolder, l logger.Logger, d string) *CreateCommand {
	create := &CreateCommand{
		blogName:       "",
		blogScaffolder: s,
		logWriter:      l,
		targetDir:      "",
	}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Teel blog",
		Run:   create.runCreate,
	}

	cmd.Flags().StringVarP(&create.blogName, "name", "n", "", "Name of blog to create")
	cmd.Flags().StringVarP(&create.targetDir, "target", "t", d, "Target directory to create blog")
	create.Command = cmd

	return create
}

func (c *CreateCommand) runCreate(cmd *cobra.Command, a []string) {
	c.logWriter.Log("--create command invoked")
	fmt.Printf("Creating blog server at %s\n", c.targetDir)

	scaffoldErr := c.blogScaffolder.Scaffold(c.targetDir, c.blogName)

	if scaffoldErr != nil {
		c.logWriter.Errorf("--create command errored with: %v\n", scaffoldErr)
		fmt.Printf("Error: failed to create blog server: %s\n", scaffoldErr.Error())
		c.logWriter.Log("--create command failed")

		return
	}

	fmt.Println("Blog server created successfully")
	c.logWriter.Log("--create command completed")
}
