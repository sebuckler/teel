package commands

import (
	"fmt"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/spf13/cobra"
)

type CreateCommand struct {
	blogName   string
	scaffolder scaffolder.Scaffolder
	logger     logger.Logger
	targetDir  string
	*cobra.Command
}

func NewCreate(s scaffolder.Scaffolder, l logger.Logger, d string) *CreateCommand {
	create := &CreateCommand{
		blogName:   "",
		scaffolder: s,
		logger:     l,
		targetDir:  "",
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

func (c *CreateCommand) runCreate(*cobra.Command, []string) {
	c.logger.Log("--create command invoked")
	fmt.Printf("Creating blog server at %s\n", c.targetDir)

	scaffoldErr := c.scaffolder.Scaffold(c.targetDir, c.blogName)

	if scaffoldErr != nil {
		c.logger.Errorf("--create command errored with: %v\n", scaffoldErr)
		fmt.Printf("Error: failed to create blog server: %s\n", scaffoldErr.Error())
		c.logger.Log("--create command failed")

		return
	}

	fmt.Println("Blog server created successfully")
	c.logger.Log("--create command completed")
}
