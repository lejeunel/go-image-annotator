{
  perSystem =
    { pkgs, ... }:
    {
      packages.nix-cache-sync = pkgs.writeShellApplication {
        name = "nix-cache-sync";
        runtimeInputs = [
          pkgs.findutils
          pkgs.coreutils
        ];
        text = builtins.readFile ./../nix-cache-sync.sh;
      };
    };

}
