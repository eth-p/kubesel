# kubesel

>[!important]
> **This is a major work in progress.**
>
> The code is messy/undocumented, the UX needs improving, and
> there's a lot of things that still need to be implemented.

Kubesel (**kube**config **sel**ector) is your modern approach to working with
[kubectl](https://kubernetes.io/docs/reference/kubectl/) configuration in a
multi-cluster, multi-namespace environment. Quickly and easily change your
active kubectl context, namespace, and cluster through a single program.

Designed from the ground up using the [KUBECONFIG environment variable](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable),
kubesel creates and manages a unique kubeconfig file for each shell session.
If you change your cluster in one pane, there's no chance\* of accidentally
running `kubectl delete` on that cluster from your other pane.

\*unless you copied the `KUBECONFIG` environment variable to the other pane.

## 🚧 Progress Checklist 🚧

**Features:**

 - [x] Switch contexts with `kubesel context <context>`
 - [x] Switch users with `kubesel user <user>`
 - [x] Switch clusters with `kubesel cluster <cluster>`
 - [x] Switch namespaces with `kubesel namespace <namespace>`
 - [ ] Listing namespaces with `kubesel list namespaces`
 - [ ] Automatically switch users with `kubesel cluster <clusters>`
 - [ ] Garbage collection of outdated/defunct session files
 - [ ] Automatic garbage collection
 - [ ] `kubesel status` to show current session info
 - [ ] `kubesel status` to show problems with kubeconfig files

**User Experience:**

 - [x] Use fzf-based UI to select context/user/cluster when none provided.
 - [x] Use fzf-based to fuzzily match clusters, showing a UI for multiple matches.
 - [x] Init script for `fish`
 - [x] Init script for `bash`
 - [x] Init script for `zsh`
 - [ ] Shell completions for `fish`
 - [ ] Shell completions for `zsh`
 - [ ] Better error handling
 - [ ] Terminal colors
 - [ ] Automatic detection of terminal color support
 - [ ] Consistent exit codes
 - [ ] Install from nix flakes
 - [ ] Install from GitHub releases
 - [ ] Safety kubeconfig file
 - [x] Manpages

**Developer Experience:**

 - [ ] Set up devenv
 - [ ] Set up golangci-lint

## Installation

**With Go:**

```bash
go install github.com/eth-p/kubesel/cmd/kubesel
```

### Setup

> [!important]
> The `KUBECONFIG` environment variable should be set before kubesel is run.

In order for kubesel to manage a per-shell kubeconfig file, it needs to create
the file and update your `KUBECONFIG` environment variable when the shell
starts.

#### Fish

Add this to your `~/.config/fish/config.fish` file:

```fish
kubesel init fish | source
```

## Usage

**Change Cluster, User, or Namespace:**
```bash
kubesel cluster my-cluster      # use this cluster
kubesel user my-user            # use this user
kubesel namespace my-namespace  # use this namespace
```

**Change Cluster, User, _and_ Namespace:**
```bash
# Use the cluster, user, and namespace from this context.
kubesel context my-context
```

**View Contexts, Clusters, or Users:**
```bash
kubesel list clusters
kubesel list contexts
kubesel list users
```

## Tips

### Alternate kubsel list outputs

The `kubesel list` command supports changing its output format with `--output`.  
Supported formats are:

 - `list` for just the names in an unsorted list
 - `table` for a table
 - `col` for columns
 - `col=COL1,COL2` for specific columns


## Alternatives

### kubectx
https://github.com/ahmetb/kubectx

 - ✅ Shell completions.
 - ✅ Fuzzy-finding.
 - ⚠️ Changes affect all shells.

### kubesess
https://github.com/Ramilito/kubesess

 - ✅ Per-shell cluster/namespace/context.
 - ✅ Fuzzy-finding.
 - ⚠️ Does not handle OIDC refresh tokens properly.

### fish-kubeswitch
https://github.com/eth-p/fish-kubeswitch

 - ✅ Per-shell cluster/namespace/context.
 - ✅ Shell completions.
 - ⚠️ Only supports [fish shell](https://fishshell.com/).
 - ⚠️ Wraps kubectl as a shell function.
