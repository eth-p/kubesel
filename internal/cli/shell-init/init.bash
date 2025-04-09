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
    if type -t _get_comp_words_by_ref &>/dev/null && complete -p @@KUBESEL_BASENAME@@ &>/dev/null; then
        source <(@@KUBESEL@@ completion bash)
    fi
}

__kubesel_load_completions
