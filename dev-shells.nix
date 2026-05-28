{
  perSystem =
    { pkgs, ... }:
    {

      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          gotestsum
          oapi-codegen
          tailwindcss_4
        ];
      };
    };

}
