{
  description = "A modern approach to managing kubectl in multi-cluster environments.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    systems.url = "github:nix-systems/default";
    gomod2nix = {
      url = "github:nix-community/gomod2nix?rev=9b6058eb92b71f2aa20c86a118aab9e41632e8bc";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      systems,
      gomod2nix,
    }:
    let
      rev = self.rev or self.dirtyRev or "unknown";
      forEachSystem =
        fn:
        nixpkgs.lib.genAttrs (import systems) (
          system:
          fn {
            inherit system;
            pkgs = (import nixpkgs { inherit system; });
          }
        );
    in
    {
      packages = forEachSystem (
        { system, pkgs }:
        import ./nix/packages {
          inherit pkgs;
          inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
          flakeGitRev = rev;
          flakeRoot = ./.;
        }
      );

      devShells = forEachSystem (
        { system, pkgs }:
        {
          default = pkgs.mkShell {
            packages = [
              gomod2nix.legacyPackages.${system}.gomod2nix
            ];
          };
        }
      );

      overlays.default = final: prev: {
        kubesel = self.packages.${prev.stdenv.hostPlatform.system}.kubesel;
      };
    };
}
