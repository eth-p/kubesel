{
  pkgs,

  # from flake.nix
  flakeGitRev,
  flakeRoot,

  # from gomod2nix
  buildGoApplication,
}:
rec {
  default = kubesel;
  kubesel = pkgs.callPackage ./kubesel.nix {
    inherit flakeGitRev flakeRoot buildGoApplication;
  };
}
