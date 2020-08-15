package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/x-color/store/store"
)

func runAddCmd(cmd *cobra.Command, args []string) error {
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

	err = store.Add(data, keys, args[1])
	if err != nil {
		return err
	}

	return store.Save(data, f)
}

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <key> <value>",
		Short:   "Add data",
		Example: "  store add .path.to value",
		Args:    cobra.ExactArgs(2),
		RunE:    runAddCmd,
	}

	return cmd
}
