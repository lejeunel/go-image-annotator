{
  perSystem =
    { pkgs, config, ... }:
    {
      packages.datahub =
        let
          gomdxp = config.packages.gomdxp;
        in

        pkgs.buildGoModule {
          name = "datahub";
          src = ./.;
          vendorHash = "sha256-dhNo4P2xKs5tKpwgFTX96sN8d1T+JWYOiqhOrci4qnM=";
          buildInputs = with pkgs; [
            git
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
            echo "Building documentation"
            ${gomdxp}/bin/docexport compile ./assets/docs ./site/docs

            echo "Installing tailwind plugins"
            mkdir -p tmp/node_modules
            npm install --prefix tmp @tailwindcss/typography

            echo "Generating CSS"
            ${pkgs.tailwindcss}/bin/tailwindcss \
                -i app/app.css \
                -o site/static/styles.css \
                -c tailwind.config.js
          '';

          buildPhase = ''
            echo "Building datahub..."
            go build -o $out/bin/datahub -ldflags "-s -w" ./cmd
          '';
        };

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
