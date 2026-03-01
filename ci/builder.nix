{
  perSystem =
    { pkgs, config, ... }:
    {
      packages.builder = pkgs.dockerTools.buildImage {
        name = "datahub-builder";
        tag = "latest";
        created = "now";

        copyToRoot = pkgs.buildEnv {
          name = "image-root";
          pathsToLink = [
            "/bin"
            "/etc"
            "/tmp"
            "/var"
          ];
          paths = with pkgs; [
            config.packages.nix-cache-sync
            bashInteractive
            cacert
            coreutils
            git
            skopeo
            nix
            gnugrep
            (fakeNss.override {
              extraPasswdLines = [
                "nixbld1:x:997:996:Nix build user 1:/var/empty:/usr/sbin/nologin"
                "nobody:x:65534:65524:nobody:/var/empty:/bin/sh"
              ];
              extraGroupLines = [
                "nixbld:x:996:nixbld1"
                "nobody:x:65534:"
              ];
            })
            (writeTextDir "etc/nix/nix.conf" ''
              sandbox = false
              experimental-features = nix-command flakes
            '')
            (writeTextDir "etc/containers/policy.json" ''
              { "default" : [ { "type": "insecureAcceptAnything" } ] }
            '')
            (runCommand "tmp" { } "mkdir -p $out/tmp $out/var/tmp")
            dockerTools.caCertificates
          ];
        };

        config = {
          Cmd = [ "${pkgs.bashInteractive}/bin/bash" ];
          Env = [
            "NIX_PAGER=cat"
            "USER=nobody"
          ];
        };
      };
    };

}
