{
  description = "A modern approach to managing kubectl in multi-cluster environments.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      gomod2nix,
    }@inputs:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        inherit (pkgs.nix-gitignore) gitignoreSource;
        inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
      in
      {
        packages = rec {
          default = kubesel;
          kubesel = buildGoApplication {
            name = "kubesel";
            pwd = ./.;
            src = gitignoreSource [ ] ./.;
            modules = ./gomod2nix.toml;

            outputs = [
              "out"
              "man"
            ];

            buildPhase = ''
              echo "compiling kubesel"
              mkdir -p $out/bin
              go build -o $out/bin/kubesel \
                -ldflags '-X main.GIT_REVISION=${self.rev or self.dirtyRev or "unknown"}' \
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
          };
        };

        devShells.default = pkgs.mkShell {
          packages = [
            gomod2nix.legacyPackages.${system}.gomod2nix
          ];
        };
      }
    ));
}
