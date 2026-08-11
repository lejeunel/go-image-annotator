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
          gotestsum
          tailwindcss_4
          oapi-codegen
          redocly
        ];
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
