{
  pkgs,
  lib,
  config,
  inputs,
  ...
}:
{
  packages = [
    pkgs.git

    # formatting
    pkgs.treefmt
    pkgs.nixfmt
    pkgs.prettier
    pkgs.tombi
  ];

  languages.go.enable = true;

  # Disable the default enterShell task.
  enterShell = lib.mkForce "";
}
