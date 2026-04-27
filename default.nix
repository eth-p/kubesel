# This file re-exports the flake.nix outputs from outside of a flake context.
# The flake.lock file is used to determine flake inputs (and their inputs).
let
  lockfile = ./flake.lock;
  lock = builtins.fromJSON (builtins.readFile lockfile);

  # resolveInput :: string -> node
  # resolveInput :: [string] -> node
  resolveInput =
    input:
    if (builtins.typeOf input) == "list" then
      # Follows (["a" "b" "c"])
      (
        let
          start = lock.nodes.${builtins.head input};
          follow = cur: next: resolveInput cur.inputs.${next};
        in
        builtins.foldl' follow start (builtins.tail input)
      )
    else
      # A node name.
      lock.nodes.${input};

  # inputNodesOf returns an attrset containing the inputs (as nodes) of the
  # provided node.
  #
  # inputNodesOf :: node -> { [string] :: node }
  inputNodesOf = node: if node ? inputs then builtins.mapAttrs (k: resolveInput) node.inputs else { };

  # loadInput returns the input's flake outputs (if it has a flake.nix) or
  # directory.
  #
  # loadInput :: node -> { [string] :: any }
  # loadInput :: node -> str
  loadInput =
    node:
    let
      input = builtins.fetchTree node.locked;
    in
    if builtins.pathExists (input + "/flake.nix") then loadFlakeInput node input else input.outPath;

  # loadFlakeInput :: node -> { [string] :: any }
  loadFlakeInput =
    node: input:
    let
      flake = import (input + "/flake.nix");
      flakeInputs = builtins.mapAttrs (k: loadInput) (inputNodesOf node);
      outputs =
        flake.outputs (
          flakeInputs
          // {
            self = outputs;
          }
        )
        // {
          __toString = _: input.outPath;
        };
    in
    outputs;
in
loadFlakeInput lock.nodes.root ./.
