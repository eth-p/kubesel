function __kubesel_init
    set -l new_kubeconfig (@@KUBESEL@@ __init --pid=$fish_pid)
    if test $status -eq 0
        set -x KUBECONFIG "$new_kubeconfig"
    end
end
__kubesel_init
functions -e __kubesel_init
