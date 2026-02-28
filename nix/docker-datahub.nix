{
  perSystem =
    { pkgs, config, ... }:
    {
      packages.docker-datahub =
        let
          datahub = config.packages.datahub;
        in

        pkgs.dockerTools.buildImage {
          name = "datahub";
          tag = "latest";
          created = "now";

          # Use Nix's built-in minimal environment
          copyToRoot = pkgs.buildEnv {
            name = "minimal-root";
            paths = [
              datahub
              pkgs.bashInteractive
              pkgs.cacert
            ];
            pathsToLink = [ "/bin" ];
            extraOutputsToInstall = [ "out" ]; # Only runtime files
          };
          config.Cmd = [ "/bin/datahub" ];
        };

    };

}
