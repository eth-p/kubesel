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
    if type -t _get_comp_words_by_ref &>/dev/null && complete -p {{ .kubesel_name | shellquote }} &>/dev/null; then
        source <({{ .kubesel_executable | shellquote }} completion bash)
    fi
}

__kubesel_load_completions
{{- end }}
