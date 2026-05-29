{
  perSystem =
    { pkgs, ... }:
    {

      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          gocyclo
          gotestsum
          oapi-codegen
          tailwindcss_4
        ];
      };
    };

}
