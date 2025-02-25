package use

import (
	"github.com/kloudlite/kl/domain/client"
	"github.com/kloudlite/kl/domain/server"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/text"
	"github.com/spf13/cobra"
)

var accCmd = &cobra.Command{
	Use:   "account",
	Short: "Switch account",
	Run: func(cmd *cobra.Command, _ []string) {
		accountName := fn.ParseStringFlag(cmd, "account")

		acc, err := server.SelectAccount(accountName)
		if err != nil {
			fn.PrintError(err)
			return
		}

		if err := client.SetAccountToMainCtx(acc.Metadata.Name); err != nil {
			fn.PrintError(err)
			return
		}

		fn.Logf("%s %s", text.Blue(text.Bold("\nSelected Account:")), acc.Metadata.Name)
	},
}

func init() {
	accCmd.Flags().StringP("account", "a", "", "account name")
	accCmd.Aliases = append(accCmd.Aliases, "acc")
}
