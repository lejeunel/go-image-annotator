{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };
  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {

      systems = [ "x86_64-linux" ];

      imports = [
        ./nix/datahub.nix
        ./nix/docker-datahub.nix
        ./nix/gomdxp.nix
        ./nix/devshells.nix
        ./nix/nix-cache-sync.nix
        ./nix/tester.nix
        ./nix/builder.nix
      ];
    };
}
