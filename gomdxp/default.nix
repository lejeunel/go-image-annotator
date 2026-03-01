{
  perSystem =
    { pkgs, ... }:
    {
      packages.gomdxp = pkgs.buildGoModule {
        name = "gomdxp";
        src = ./.;
        vendorHash = "sha256-R6vXs+Tkws1Oj1DcH945mUG/mm0gE0JzSVbltxVXSL4=";
        buildInputs = with pkgs; [
          go
        ];
        env.CGO_ENABLED = 0;
        trimPath = true;
        ldflags = [
          "-s" # Omit symbol table
          "-w" # Omit DWARF symbols
        ];
      };

    };

}
