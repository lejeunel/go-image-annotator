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
          tailwindcss
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
