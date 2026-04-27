{
  lib,
  nix-gitignore,

  # from flake.nix (via root default.nix)
  flakeGitRev,
  flakeRoot,

  # from gomod2nix
  buildGoApplication,
}:
buildGoApplication {
  name = "kubesel";
  pwd = flakeRoot;
  src = nix-gitignore.gitignoreSource [ ] flakeRoot;
  modules = flakeRoot + "/gomod2nix.toml";

  outputs = [
    "out"
    "man"
  ];

  buildPhase = ''
    echo "compiling kubesel"
    mkdir -p $out/bin
    go build -o $out/bin/kubesel \
      -ldflags '-X main.GIT_REVISION=${flakeGitRev}' \
      -tags no_init_completions \
      ./

    echo "generating fish completions"
    mkdir -p $out/share/fish/vendor_completions.d
    $out/bin/kubesel completion fish > $out/share/fish/vendor_completions.d/kubesel.fish

    echo "generating zsh completions"
    mkdir -p $out/share/zsh/site-functions
    $out/bin/kubesel completion zsh > $out/share/zsh/site-functions/_kubesel

    echo "generating bash completions"
    mkdir -p $out/share/bash-completion
    $out/bin/kubesel completion bash > $out/share/bash-completion/kubesel.bash

    echo "generating manuals"
    mkdir -p $man/share/man/man1
    go run hack/generate-man.go -outdir $man/share/man/man1
  '';

  meta = {
    description = "A modern approach to managing kubectl in multi-cluster environments.";
    license = lib.licenses.mit;
    mainProgram = "kubesel";
  };
}
