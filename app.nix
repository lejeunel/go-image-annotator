{ ... }:
let
  version = "0.1.0";
in
{
  perSystem =
    {
      pkgs,
      ...
    }:
    let
      app = pkgs.buildGoModule {
        pname = "go-image-annotator";
        inherit version;

        src = ./.;

        vendorHash = "sha256-rgTdLVCVA5QopmuIC9H4OPTG6pgoXcG1t2xNH1dRVhA";

        nativeBuildInputs = with pkgs; [
          go
          git
          gnumake
        ];

        subPackages = [ "." ];

        ldflags = [
          "-X github.com/lejeunel/go-image-annotator/globals.Version=${version}"
          "-X 'github.com/lejeunel/go-image-annotator/globals.Date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
        ];
      };
      
      tests = app.overrideAttrs (old: {
        pname = "${old.pname}-tests";

        # Run `go test` as the check phase instead of building the binary.
        doCheck = true;
        checkFlags = [ "-v" ];
        checkPhase = ''
              runHook preCheck
              go test -v -race ./...
              runHook postCheck
              '';
          
        # Skip the actual binary build+install; we only want the test run.
        dontBuild = true;
        installphase = ''
          mkdir -p $out
          touch $out/tests-passed
        '';
      });
      
      format = pkgs.stdenvNoCC.mkDerivation {
        pname = "go-image-annotator-format-check";
        inherit version;

        src = ./.;

        nativeBuildInputs = with pkgs; [
          gofumpt
          golines
        ];

        dontBuild = true;
        doCheck = true;

        checkPhase = ''
          runHook preCheck

          unformatted=$(gofumpt -l .)
          if [ -n "$unformatted" ]; then
            echo "gofumpt: the following files are not formatted:"
            echo "$unformatted"
            exit 1
          fi

          toolong=$(golines -l .)
          if [ -n "$toolong" ]; then
            echo "golines: the following files need reformatting:"
            echo "$toolong"
            exit 1
          fi

          runHook postCheck
        '';

        installPhase = ''
          mkdir -p $out
          touch $out/format-ok
        '';
      };
    in
    {
      packages.default = app;
      packages.app = app;

      checks.default = tests;
      checks.tests = tests;

      checks.format = format;
    };
}
