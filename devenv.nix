{
  pkgs,
  lib,
  config,
  inputs,
  ...
}:
let
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in
{
  packages = [
    pkgs.git

    # formatting
    pkgs.treefmt
    pkgs.nixfmt-rfc-style
    pkgs.nodePackages.prettier
    pkgs.nodePackages.prettier-plugin-toml
  ];

  languages.go.enable = true;
  languages.go.package = pkgs-unstable.go;

  # Formatting.
  scripts.prettier.exec = ''
    ${pkgs.nodePackages.prettier}/bin/prettier \
      --plugin ${pkgs.nodePackages.prettier-plugin-toml}/lib/node_modules/prettier-plugin-toml/lib/index.js \
      "$@" || exit $?
  '';

  # Disable the default enterShell task.
  enterShell = lib.mkForce "";
}
