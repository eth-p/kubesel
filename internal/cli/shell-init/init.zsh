# Update the KUBECONFIG environment variable.
__kubesel_init() {
    unset -f __kubesel_init

    {{- with .add_kubeconfigs }}
    export KUBECONFIG="$KUBECONFIG:"{{ join . ":" | shellquote }}
    {{- end }}

    local new_kubeconfig
    new_kubeconfig="$({{ .kubesel_executable | shellquote }} __init --pid=$$)"
    if test $? -eq 0; then
        export KUBECONFIG="$new_kubeconfig"
    fi
}

__kubesel_init

{{- if .load_completions }}
# Load completions.
__kubesel_load_completions() {
    unset -f __kubesel_load_completions
    local kubesel_name={{ .kubesel_name | shellquote }}
    if type compdef &>/dev/null && test -z "$_comps[$kubesel_name]"; then
        source <({{ .kubesel_executable | shellquote }} completion zsh)
    fi
}

__kubesel_load_completions
{{- end }}
