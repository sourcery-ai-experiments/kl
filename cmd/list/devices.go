package list

import (
	"errors"
	"fmt"
	"github.com/kloudlite/kl/domain/client"
	server2 "github.com/kloudlite/kl/domain/server"
	common_util "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/table"
	"github.com/kloudlite/kl/pkg/ui/text"
	"strings"

	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "list all your devices in selected account",
	Long: `List all your devices in selected account.
Examples:
	# list all the devices with selected account
  kl list devices

	# list all the devices with accountId
  kl list devices <accountId>
`,
	Run: func(_ *cobra.Command, args []string) {
		err := Devices(args)
		if err != nil {
			common_util.PrintError(err)
			return
		}
	},
}

func Devices(args []string) error {

	var devices []server2.Device
	var err error

	rs, err := server2.GetRegions()
	if err != nil {
		return err
	}

	getRegionName := func(regionId string) string {
		for _, r2 := range rs {
			if r2.Region == regionId {
				return r2.Name
			}
		}

		return regionId
	}

	if len(args) >= 1 {
		devices, err = server2.GetDevices(common_util.MakeOption("accountId", ""))
	} else {
		devices, err = server2.GetDevices()
	}

	if err != nil {
		return err
	}

	if len(devices) == 0 {
		return errors.New("no devices found")
	}

	header := table.Row{
		table.HeaderText("name, id"),
		table.HeaderText("Active Region"),
		table.HeaderText("exposed ports"),
	}

	cDid, _ := server2.CurrentDeviceId()

	rows := make([]table.Row, 0)

	for _, a := range devices {
		rows = append(rows, table.Row{

			func() string {
				if cDid == a.Id {
					return text.Colored(fmt.Sprintf("*%s, %s", a.Name, a.Id), 2)
				}
				return fmt.Sprintf("%s, %s", a.Name, a.Id)
			}(),

			func() string {
				if cDid == a.Id {
					return fmt.Sprintf("%s\n%s", text.Colored(a.Region, 2), text.Colored(getRegionName(a.Region), 2))
				}
				return a.Region
			}(),

			strings.Join(func() []string {
				var ports []string
				for i, v := range a.Ports {
					prt := fmt.Sprintf("%s%d:%d", func() string {
						if (i+1)%3 == 0 {

							if cDid == a.Id {
								return fmt.Sprint("\n", text.Color(2))
							}
							return "\n"
						}
						return ""
					}(), v.Port, func() int {
						if v.TargetPort == 0 {
							return v.Port

						}
						return v.TargetPort
					}())

					ports = append(ports, func() string {
						if cDid == a.Id {
							return text.Colored(prt, 2)
						}
						return prt
					}())
				}
				return ports
			}(), func() string {
				if cDid == a.Id {
					return text.Colored(", ", 2)
				}
				return ", "
			}()),
		})
	}

	fmt.Println(table.Table(&header, rows))

	if accountId, _ := client.CurrentAccountName(); accountId != "" {
		table.KVOutput("devices of", accountId, true)
	}

	table.TotalResults(len(devices), true)
	return nil
}

func init() {
	devicesCmd.Aliases = append(devicesCmd.Aliases, "device")
	devicesCmd.Aliases = append(devicesCmd.Aliases, "dev")
}
