<img src="assets/kubectl-epinio.png" width="480px" alt="kubectl-epinio logo"/>


# kubectl-epinio

[![CI](https://github.com/enrichman/kubectl-epinio/actions/workflows/main.yml/badge.svg)](https://github.com/enrichman/kubectl-epinio/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/enrichman/kubectl-epinio)](https://goreportcard.com/report/github.com/enrichman/kubectl-epinio)

With this kubectl plugin you will be able to run some administrative command for [Epinio](https://github.com/epinio/epinio) in a convenient way.

For example you can create, describe and edit Epinio users and roles, with the help of autocompletions, and within an interactive mode!

```
-> % kubectl epinio create user myuser -i
Password: 
Retype password: 
? Namespaces assigned to the user: workspace
? Global Roles assigned to the user: user, my-role
Do you want to assign Namescoped Roles? [y/n] y
? Namescoped Roles assigned to the user: reader2, reader
? Namespaces for 'reader2' role: workspace
? Namespaces for 'reader' role: workspace

Username:       myuser
Password:       $2a$10$CI9Vg03zJ1rR6jpTHwj62OS0r/.LtJ/u8h0A/TTsfO.64LizCqfiy
Roles:          user
                my-role
                reader2:workspace
                reader:workspace

Namespaces:     workspace


Create? [y/n] y
User created!
```

# Installation

To install the plugin you can download a binary from the release page:

```shell
wget -qO- https://github.com/enrichman/kubectl-epinio/releases/download/v0.0.1/kubectl-epinio_linux_x86_64.tar.gz | \
    sudo tar -xvz -C /usr/local/bin kubectl-epinio
```

### Autocompletion

To enable the [autocompletion](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#enable-shell-autocompletion) for the plugins you need to have at least the v1.26 of `kubectl` (see [#105867](https://github.com/kubernetes/kubernetes/pull/105867) for more info about this).

You can use this plugin to generate it: https://github.com/marckhouzam/kubectl-plugin_completion

```
kubectl krew install --manifest-url https://raw.githubusercontent.com/marckhouzam/kubectl-plugin_completion/v0.1.0/plugin-completion.yaml
kubectl plugin-completion generate
export PATH=$PATH:$HOME/.kubectl-plugin-completion
```

## Testing

Run the following:
1. `make infra-setup`
2. `make test`
3. (optional) `make infra-teardown`
