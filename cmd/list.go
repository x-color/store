package cmd

import (
	"github.com/spf13/cobra"
	"github.com/x-color/store/store"
)

func runListCmd(cmd *cobra.Command, args []string) error {
	var keys []string
	var err error
	if len(args) > 0 {
		keys, err = store.SplitKeys(args[0])
		if err != nil {
			return err
		}
	}

	f, err := store.DataFile()
	if err != nil {
		return err
	}

	data, err := store.Load(f)
	if err != nil {
		return err
	}

	keyValues, err := store.Get(data, keys)
	if err != nil {
		return err
	}

	for k := range keyValues {
		cmd.Println(k)
	}

	return nil
}

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <key>",
		Short:   "List keys",
		Example: "  store list .path.to",
		Args:    cobra.MaximumNArgs(1),
		RunE:    runListCmd,
	}

	return cmd
}
