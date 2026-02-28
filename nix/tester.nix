{
  perSystem =
    { pkgs, ... }:
    {
      packages.tester = pkgs.dockerTools.buildImage {

        created = "now";
        name = "datahub-tester";
        tag = "latest";
        copyToRoot = pkgs.buildEnv {
          name = "image-root";
          pathsToLink = [
            "/bin"
            "/etc"
            "/tmp"
            "/var"
          ];
          paths = with pkgs; [
            bashInteractive
            cacert
            coreutils
            git
            dockerTools.caCertificates
            go
            gcc14
            gnugrep
          ];
        };

        config = {
          Cmd = [ "${pkgs.bashInteractive}/bin/bash" ];
        };
      };

    };

}
