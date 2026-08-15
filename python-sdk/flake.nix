{
  description = "Python environment with uv and pyright";

  inputs = {
    flake-compat.url = "github:edolstra/flake-compat";
    flake-compat.flake = false;
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs =
    {
      flake-utils,
      nixpkgs,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            allowUnfree = true;
          };
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            uv
            python312
            ruff
            basedpyright

          ];

          shellHook = ''
            export LD_LIBRARY_PATH="${pkgs.glib.out}/lib:${pkgs.stdenv.cc.cc.lib}/lib:${pkgs.zlib}/lib:${pkgs.libGL}/lib:$LD_LIBRARY_PATH"
          '';
        };
      }
    );
}
