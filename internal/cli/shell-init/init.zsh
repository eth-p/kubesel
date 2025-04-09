# Update the KUBECONFIG environment variable.
__kubesel_init() {
    unset -f __kubesel_init

    local new_kubeconfig
    new_kubeconfig="$(@@KUBESEL@@ __init --pid=$$)"
    if test $? -eq 0; then
        export KUBECONFIG="$new_kubeconfig"
    fi
}

__kubesel_init

# Load completions.
__kubesel_load_completions() {
    unset -f __kubesel_load_completions
    if type compdef &>/dev/null && test -z "$_comps[@@KUBESEL_BASENAME@@]"; then
        source <(@@KUBESEL@@ completion zsh)
    fi
}

__kubesel_load_completions
