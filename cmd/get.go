package cmd

import (
	"github.com/spf13/cobra"
	"github.com/x-color/store/store"
)

func runGetCmd(cmd *cobra.Command, args []string) error {
	keys, err := store.SplitKeys(args[0])
	if err != nil {
		return err
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

	if len(keyValues) == 1 {
		if v, ok := keyValues["."]; ok {
			cmd.Println(v)
			return nil
		}
	}

	for k, v := range keyValues {
		cmd.Printf("%v = %v\n", k, v)
	}

	return nil
}

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <key>",
		Short:   "Get data",
		Example: "  store get .path.to",
		Args:    cobra.ExactArgs(1),
		RunE:    runGetCmd,
	}

	return cmd
}
