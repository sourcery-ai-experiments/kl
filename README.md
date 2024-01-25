# Kloudlite CLI `kl`

This CLI help you to work with kloudlite using your terminal.

### Installation

To install the latest version in Linux or Mac you can run the following command on terminal.

##### Install latest with curl
```sh
curl https://i.jpillora.com/kloudlite/kl! | bash
```

##### Install latest with wget

```sh
wget -qO- https://i.jpillora.com/kloudlite/kl! | bash
```

##### Install specific version
```sh
curl https://i.jpillora.com/kloudlite/kl@v1.0.0! | bash
```

##### download but don't install
```sh
curl https://i.jpillora.com/kloudlite/kl | bash
```

### Authentication

To login and logout you can use the following commands.

```sh
kl auth login
kl auth logout
kl auth status
```

### Select Account

To select account you can use the following command.

```sh
kl switch account
```

### Initialize your workspace
To work with any project you need to initialize your workspace where you can define 
environments, managed resouces, mounts and etc.
To initialize you workspace you can use the following command.
```sh
kl init
```


### Listing Resources

With this CLI you can list accounts, projects, envs, devices, configs, secrets and apps.
To list resources you can use the following commands.
For more details visit [kl list](./docs/kl_list.md)


```sh
kl list accounts
kl list projects
kl list envs
kl list devices
kl list configs
kl list secrets
kl list apps
```

with these commands you can provide the resource id. In case of you don't provide resource 
it it will show you a picker. For more details visit [kl list](./docs/kl_list.md)

### Working with vpn

To access services of cluster and tunnel your local app to the server you need to connect to vpn.
For that you can use the following commands.

```sh
sudo kl vpn start
sudo kl vpn stop
sudo kl vpn status

# to tunnel traffic to your local you need to expose ports also
kl wg expose -p <server_port>:<local_port>
kl wg expose -p <server_port>:<local_port> -d    # provide -d flag to delete
```

### Working with environments
We support multiple environments to work with. these commands 
will help you to import config,secrets as environment variables, mount and also ipmprting managed resources.

For more details visit [kl add](./docs/kl_add.md).


```sh
# Adding
kl add config
kl add secret
kl add mres
kl add mount <file_path/file_name>
```

### Intercepting App
You can tunnel you local running app to the server and intercept your app to forward all the request of that app to your local system. 
for that you need to perform following actions.
- [connecting to vpn](./docs/kl_vpn_start.md)
- [exposing port](./docs/kl_vpn_expose.md)
- [intercept an app](./docs/kl_vpn_intercept.md)

So you can use following commands to work with interception. 
For more details visit [kl vpn intecept start](./docs/kl_vpn_intercept_start.md) and [kl vpn intecept stop](./docs/kl_vpn_intercept_stop.md)

```sh
kl vpn intercept start
kl vpn intercept stop
```

### KL Config File structure
This is the structure of app config file which will be generated by executing the command `kl -- <cmd>` and 
you can also modify this file according to your requirement.
```yaml
version: v1
name: <project_name>
mres: 
- name: service/<mres_name>
  env:
  - name: <env_name>
    key: <local_key>
    refkey: <server_key>
configs:
- name: <config_name>
  env:
  - key: <local_key>
    refkey: <server_key>
secrets:
- name: <secret_name>
  env:
  - key: <local_key> 
    refkey: <server_key> 
env:
- key: <env_key>    # eg. NODE_ENV
  value: <env_value>    # eg. development
fileMount:
  mountBasePath: <base_mount_path> # eg. ./.mounts
  mounts:
  - path: <mount_path> # eg. /tmp
    type: <type> # eg. config or secret
    name: <config_name>
```

## Getting All environments 
According to above config file you can get all the environments to your local shell.
you can use the following commands for getting all environments to working shell.
```
kl -- <cmd>  # will execute the command with all the environments eg. kl -- npm start
kl -- <shell>  # will open the shell with all the environments eg. kl -- bash

kl -- printenv # will print all the environments
```



> This CLI is under development so, more information will will be updated in this doc. also if some new commands will be added to the CLI will be updated to this doc.


> for more details visit [docs](./docs/kl.md)
