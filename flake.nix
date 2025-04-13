{
  description = "A modern approach to managing kubectl in multi-cluster environments.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
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

            outputs = [ "out" "man" ];

            buildPhase = ''
              echo "compiling kubesel"
              mkdir -p $out/bin
              go build -o $out/bin/kubesel ./cmd/kubesel

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
