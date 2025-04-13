{ pkgs, lib, config, inputs, ... }:
let
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in
{
  packages = [
    pkgs.git
    pkgs.goreleaser
  ];

  languages.go.enable = true;
  languages.go.package = pkgs-unstable.go;

  # Disable the default enterShell task.
  enterShell = lib.mkForce "";
}
