__kubesel_init() {
    local new_kubeconfig
    new_kubeconfig="$(@@KUBESEL@@ __init --pid=$$)"
    if test $? -eq 0; then
        export KUBECONFIG="$new_kubeconfig"
    fi
}
__kubesel_init
unset -f __kubesel_init
