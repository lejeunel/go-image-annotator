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
