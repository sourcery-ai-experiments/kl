package list

import (
	"errors"
	"fmt"

	"github.com/kloudlite/kl/domain/client"
	"github.com/kloudlite/kl/domain/server"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/table"

	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "envs",
	Short: "Get list of environments in current project",
	Run: func(cmd *cobra.Command, args []string) {
		err := listEnvironments(cmd, args)
		if err != nil {
			fn.PrintError(err)
			return
		}
	},
}

func listEnvironments(cmd *cobra.Command, args []string) error {

	var err error
	envs, err := server.ListEnvs()
	if err != nil {
		return err
	}

	if len(envs) == 0 {
		return errors.New("no environments found")
	}

	env, _ := client.CurrentEnv()
	envName := ""
	if env != nil {
		envName = env.Name
	}

	header := table.Row{table.HeaderText("DisplayName"), table.HeaderText("Name"), table.HeaderText("ready")}
	rows := make([]table.Row, 0)

	for _, a := range envs {
		rows = append(rows, table.Row{
			fn.GetPrintRow(a, envName, a.DisplayName, true),
			fn.GetPrintRow(a, envName, a.Metadata.Name),
			fn.GetPrintRow(a, envName, a.Status.IsReady),
		})
	}

	fmt.Println(table.Table(&header, rows))

	if s := fn.ParseStringFlag(cmd, "output"); s == "table" {
		table.TotalResults(len(envs), true)
	}
	table.TotalResults(len(envs), true)
	return nil
}

func init() {
	fn.WithOutputVariant(envCmd)
	envCmd.Aliases = append(envCmd.Aliases, "env")
}
