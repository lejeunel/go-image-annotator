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

        vendorHash = "sha256-2rdqCI9dWug5j68jj62RrEcorz42ibEFE/7yMPv2Peg=";

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
    in
    {
      packages.default = app;
      packages.app = app;
    };
}
