{
  perSystem =
    {
      pkgs,
      ...
    }:
    {
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gnumake
          gopls
          gofumpt
          golines
          tailwindcss_4
          oapi-codegen
          redocly
          air
          lefthook
          hugo
        ];

        shellHook = ''
          if command -v lefthook &> /dev/null; then
            lefthook install
          fi
          '';
      };
      devShells.test = pkgs.mkShell {
        packages = with pkgs; [
          go
          gotestsum
        ];
      };
      devShells.format = pkgs.mkShell {
        packages = with pkgs; [
          gofumpt
          golines
        ];
      };
    };
}
