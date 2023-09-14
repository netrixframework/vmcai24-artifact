package cmd

import (
	"errors"
	"fmt"

	"github.com/netrixframework/bftsmart-testing/client"
	"github.com/spf13/cobra"
)

func ClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "client",
	}
	cmd.AddCommand(getCmd)
	cmd.AddCommand(setCmd)
	cmd.AddCommand(deleteCmd)
	return cmd
}

var getCmd = &cobra.Command{
	Use: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid number of arguments")
		}
		client := client.NewBFTSmartClient(&client.BFTSmartClientConfig{
			CodePath: "/netrixframework/bft-smart",
		})
		result, err := client.Get(args[0])
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var setCmd = &cobra.Command{
	Use: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("invalid number of arguments")
		}
		client := client.NewBFTSmartClient(&client.BFTSmartClientConfig{
			CodePath: "/netrixframework/bft-smart",
		})
		result, err := client.Set(args[0], args[1])
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid number of arguments")
		}
		client := client.NewBFTSmartClient(&client.BFTSmartClientConfig{
			CodePath: "/netrixframework/bft-smart",
		})
		result, err := client.Remove(args[0])
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}
