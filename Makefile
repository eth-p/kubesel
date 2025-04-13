# Enable ONESHELL
# https://www.gnu.org/software/make/manual/html_node/One-Shell.html
ifneq ($(filter oneshell,$(.FEATURES)),oneshell)
    $(error A newer version of Make is required)
endif

.ONESHELL:
.SILENT:

# Use devenv if it's installed and not currently active.
devenv_installed := $(if $(shell command -v devenv 2>/dev/null),true)
devenv_active    := $(if ${DEVENV_PROFILE},true)
ifeq (${devenv_installed}/${devenv_active},true/)
    SHELL := devenv
    .SHELLFLAGS := shell --quiet -- bash -e -x -c
else
	.SHELLFLAGS := -x -c
endif

## bin: compile the kubesel executable
.PHONY: bin
bin:
	go build -o "kubesel" ./cmd/kubesel

## format: reformat source code
.PHONY: format
format:
	go fmt \
		./cmd/... \
		./internal/... \
		./pkg/...

## doc: generate manpages
.PHONY: man
man:
	-$(RM) -r man/*
	mkdir -p man
	cd man && go run "../hack/generate-man.go"
