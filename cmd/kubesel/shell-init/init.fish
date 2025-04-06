function __kubesel_init
    set -l session_file (@@KUBESEL@@ __init --pid=$fish_pid)
    if test $status -eq 0
        set -x KUBECONFIG "$KUBECONFIG:$session_file"
    end
end
__kubesel_init
functions -e __kubesel_init
