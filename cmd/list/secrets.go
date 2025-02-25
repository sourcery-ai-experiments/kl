package list

import (
	"errors"
	"fmt"

	"github.com/kloudlite/kl/domain/server"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/table"

	"github.com/spf13/cobra"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Get list of secrets in current project & selected environment",
	Run: func(cmd *cobra.Command, args []string) {

		pName := ""
		if len(args) > 1 {
			pName = args[0]
		}

		sec, err := server.ListSecrets(fn.MakeOption("projectName", pName))
		if err != nil {
			fn.PrintError(err)
			return
		}

		if err := printSecrets(cmd, sec); err != nil {
			fn.PrintError(err)
			return
		}
	},
}

func printSecrets(_ *cobra.Command, secrets []server.Secret) error {
	if len(secrets) == 0 {
		return errors.New("no secrets found")
	}

	header := table.Row{
		table.HeaderText("Display Name"),
		table.HeaderText("Name"),
		table.HeaderText("entries"),
	}

	rows := make([]table.Row, 0)

	for _, a := range secrets {
		rows = append(rows, table.Row{a.DisplayName, a.Metadata.Name, fmt.Sprintf("%d", len(a.StringData))})
	}

	fmt.Println(table.Table(&header, rows))
	table.TotalResults(len(secrets), true)
	return nil
}

func init() {
	secretsCmd.Aliases = append(secretsCmd.Aliases, "secret")
	secretsCmd.Aliases = append(secretsCmd.Aliases, "sec")
	fn.WithOutputVariant(secretsCmd)
}
