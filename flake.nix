{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages = rec {
          builder-nix-cache-sync = pkgs.writeShellApplication {
            name = "builder-nix-cache-sync";
            runtimeInputs = [
              pkgs.findutils
              pkgs.coreutils
            ];
            text = builtins.readFile ./.builder-nix-cache-sync.sh;
          };

          tester = pkgs.dockerTools.buildImage {

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

          builder = pkgs.dockerTools.buildImage {
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
                builder-nix-cache-sync
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
          default = pkgs.buildGoModule rec {
            name = "datahub";
            src = ./.;
            vendorHash = "sha256-6aJbnAp4+KJ/Fz/CghxNtt6duS+3qIwJ8Z0xyztPUxc=";
            buildInputs = with pkgs; [
              git
              hugo
              go
              gnugrep
            ];
            env.CGO_ENABLED = 0;
            trimPath = true;
            ldflags = [
              "-s" # Omit symbol table
              "-w" # Omit DWARF symbols
            ];

            preBuild = ''
              echo "=== Building documentation with Hugo ==="
              cd site/docs
              ${pkgs.hugo}/bin/hugo build
              cd ../..
            '';
          };
          dockerImage =
            let
              datahub = self.packages.${system}.default;
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
        devShells = {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              hugo
              nodejs
              gopls
              gotools
              gomodifytags
              gocode-gomod
              gotest
              age
              sqlite
              tailwindcss
            ];
          };
        };
      }
    );
}
