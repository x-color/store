package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/x-color/store/store"
)

func runSetCmd(cmd *cobra.Command, args []string) error {
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

	err = store.Set(data, keys, args[1])
	if err != nil {
		return err
	}

	return store.Save(data, f)
}

func newSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set <key> <value>",
		Short:   "Set data. If key already exists, it updates value.",
		Example: "  store set .path.to value",
		Args:    cobra.ExactArgs(2),
		RunE:    runSetCmd,
	}

	return cmd
}
