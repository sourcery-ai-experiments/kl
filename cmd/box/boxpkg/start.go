package boxpkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/kloudlite/kl/constants"
	cl "github.com/kloudlite/kl/domain/client"
	proxy "github.com/kloudlite/kl/domain/dev-proxy"
	"github.com/kloudlite/kl/domain/server"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/text"
	"github.com/kloudlite/kl/pkg/wg_vpn/wgc"
)

var containerNotStartedErr = fmt.Errorf("container not started")

func (c *client) Start() error {
	// if c.spinner.Started() {
	// 	defer c.spinner.UpdateMessage("initiating container please wait")()
	// } else {
	defer c.spinner.Start("initiating container please wait")()
	// }

	if c.verbose {
		fn.Logf("starting container in: %s", text.Blue(c.cwd))
	}

	cr, err := c.getContainer(map[string]string{
		CONT_NAME_KEY: c.containerName,
		CONT_MARK_KEY: "true",
	})
	if err != nil && err != notFoundErr {
		return err
	}

	if err == nil {
		c.spinner.Stop()
		crPath := cr.Labels[CONT_PATH_KEY]

		fn.Logf("container %s already running in %s", text.Yellow(cr.Name), text.Blue(crPath))

		// if c.cwd != crPath {
		// 	fn.Printf("do you want to stop that and start here? [Y/n]")
		// } else {
		// 	fn.Printf("do you want to restart it? [y/N]")
		// }

		// var response string
		// _, _ = fmt.Scanln(&response)
		// if c.cwd != crPath && response == "n" {
		// 	return containerNotStartedErr
		// }
		//
		// if c.cwd == crPath && response != "y" {
		// 	return containerNotStartedErr
		// }

		if err := c.Stop(); err != nil {
			return err
		}

		return c.Start()
	}

	if err := c.ensurePublicKey(); err != nil {
		return err
	}

	if err := c.ensureCacheExist(); err != nil {
		return err
	}

	envs, mmap, err := server.GetLoadMaps()
	if err != nil {
		return err
	}

	// local setup
	kConf, err := c.loadConfig(mmap, envs)
	if err != nil {
		return err
	}

	c.spinner.Stop()
	if err := cl.EnsureAppRunning(); err != nil {
		return err
	}

	c.spinner.Stop()
	d, err := server.EnsureDevice()
	if err != nil {
		return err
	}

	localEnv, err := cl.CurrentEnv()
	if err != nil {
		return err
	}

	e, err := server.GetEnv(localEnv.Name)
	if err != nil {
		return err
	}

	configuration, err := base64.StdEncoding.DecodeString(d.WireguardConfig.Value)
	if err != nil {
		return err
	}

	cfg := wgc.Config{}
	err = cfg.UnmarshalText(configuration)
	if err != nil {
		return err
	}

	c.spinner.Start()

	p, err := proxy.NewProxy(c.verbose, false)
	if err != nil {
		return err
	}

	if err := p.Start(); err != nil {
		return err
	}

	td, err := os.MkdirTemp("", "kl-tmp")
	if err != nil {
		return err
	}

	defer func() {
		os.RemoveAll(td)
	}()

	if err := func() error {
		conf, err := json.Marshal(kConf)
		if err != nil {
			return err
		}

		sshPath := path.Join(xdg.Home, ".ssh", "id_rsa.pub")

		akByte, err := os.ReadFile(sshPath)
		if err != nil {
			return err
		}

		ak := string(akByte)

		akTmpPath := path.Join(td, "authorized_keys")

		akByte, err = os.ReadFile(path.Join(xdg.Home, ".ssh", "authorized_keys"))
		if err == nil {
			ak += fmt.Sprint("\n", string(akByte))
		}

		if err := os.WriteFile(akTmpPath, []byte(ak), fs.ModePerm); err != nil {
			return err
		}

		args := []string{}

		switch runtime.GOOS {
		case constants.RuntimeWindows:
			fn.Warn("docker support inside container not implemented yet")
		default:
			args = append(args, "-v", "/var/run/docker.sock:/var/run/docker.sock:ro")
		}

		if len(cfg.DNS) > 0 {
			args = append(args, []string{
				"--dns", cfg.DNS[0].To4().String(),
				"--dns-search", fmt.Sprintf("%s.svc.%s.local", e.Spec.TargetNamespace, e.ClusterName),
			}...)
		}

		sshPort, err := cl.GetAvailablePort()
		if err != nil {
			return err
		}

		localEnv.SSHPort = sshPort

		if err := cl.SelectEnv(*localEnv); err != nil {
			return err
		}

		args = append(args, []string{
			"-v", fmt.Sprintf("%s:/tmp/ssh2/authorized_keys:ro", akTmpPath),
			"-v", "kl-home-cache:/home:rw",
			"-v", "nix-store:/nix:rw",
			"-v", fmt.Sprintf("%s:/home/kl/workspace:z", c.cwd),
			"-p", fmt.Sprintf("%d:22", sshPort),
			ImageName, "--", string(conf),
		}...)

		if err := c.runContainer(ContainerConfig{
			imageName: ImageName,
			Name:      c.containerName,
			trackLogs: true,
			labels: map[string]string{
				CONT_NAME_KEY: c.containerName,
				CONT_PATH_KEY: c.cwd,
				CONT_MARK_KEY: "true",
			},
			args: args,
		}); err != nil {
			return err
		}

		return nil
	}(); err != nil {
		return err
	}

	return nil
}
