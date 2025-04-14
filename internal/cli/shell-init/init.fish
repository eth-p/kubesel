# Update the KUBECONFIG environment variable.
function __kubesel_init
    functions -e __kubesel_init
    {{- with .add_kubeconfigs }}
    set -gx KUBECONFIG "$KUBECONFIG:"{{ join . ":" | shellquote }}
    {{- end }}
    set -l new_kubeconfig ({{ .kubesel_executable | shellquote }} __init --pid=$fish_pid)
    if test $status -eq 0
        set -gx KUBECONFIG "$new_kubeconfig"
    end
end

__kubesel_init

{{- if .load_completions }}
# Load completions.
function __kubesel_load_completions
    functions -e __kubesel_load_completions
    if test (complete -c {{ .kubesel_name | shellquote }} | count) -eq 0
        {{ .kubesel_executable | shellquote }} completion fish | source
    end
end

__kubesel_load_completions
{{- end }}
