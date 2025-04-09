# Update the KUBECONFIG environment variable.
function __kubesel_init
    functions -e __kubesel_init
    set -l new_kubeconfig (@@KUBESEL@@ __init --pid=$fish_pid)
    if test $status -eq 0
        set -x KUBECONFIG "$new_kubeconfig"
    end
end

__kubesel_init

# Load completions.
function __kubesel_load_completions
    functions -e __kubesel_load_completions
    if test (complete -c @@KUBESEL_BASENAME@@ | count) -eq 0
        @@KUBESEL@@ completion fish | source
    end
end

__kubesel_load_completions
