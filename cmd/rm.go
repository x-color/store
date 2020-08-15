package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/x-color/store/store"
)

func runRmCmd(cmd *cobra.Command, args []string) error {
	keys, err := store.SplitKeys(args[0])
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return errors.New("invalid key")
	}

	f, err := store.DataFile()
	if err != nil {
		return err
	}

	data, err := store.Load(f)
	if err != nil {
		return err
	}

	isLeaf, err := store.IsLeaf(data, keys)
	if err != nil {
		return err
	}

	if !isLeaf {
		cmd.Printf("Are you remove %v? It may include some values. [y/n]: ")
		var yes string
		_, err := fmt.Scan(&yes)
		if err != nil {
			return err
		}
		if yes != "y" && yes != "Y" {
			cmd.Println("Cancel removing")
			return nil
		}
	}

	err = store.Remove(data, keys)
	if err != nil {
		return err
	}

	return store.Save(data, f)
}

func newRmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm <key>",
		Short:   "Remove data",
		Example: "  store rm .path.to",
		Args:    cobra.ExactArgs(1),
		RunE:    runRmCmd,
	}

	return cmd
}
