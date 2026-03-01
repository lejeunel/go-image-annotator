{
  perSystem =
    { pkgs, config, ... }:
    {

      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          config.packages.gomdxp
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
