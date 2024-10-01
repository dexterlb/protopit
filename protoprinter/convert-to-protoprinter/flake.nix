{
  description = "slice and print dice";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-22.11";
  };

  outputs = { self, nixpkgs }:
    let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
      lib = nixpkgs.lib;
      used-packages = [
        pkgs.prusa-slicer
        pkgs.j2cli
        # not actually necessary, they're both in coreutils, probably pulled in by bash or something
        pkgs.coreutils
        pkgs.gnused
      ];
      convert-to-protoprinter = pkgs.writeShellScriptBin "convert-to-protoprinter" ''
        set -euo pipefail

        CONFIG=${./config.ini}
        export PATH=${lib.makeBinPath used-packages}
        export LOCALE_ARCHIVE="${pkgs.glibcLocales}/lib/locale/locale-archive";

        for f in $*
        do
          with_jinja_hack_markers="''${f}_with_jinja_hack_markers"
          prusa-slicer --load "''${CONFIG}" --export-gcode --output "''${with_jinja_hack_markers}" "''${f}"

          to_be_jinjad="''${f}_to_be_jinjad"
          cat "''${with_jinja_hack_markers}" | sed 's/;!(/{/g' | sed 's/;!)/}/g' > "''${to_be_jinjad}"
          rm "''${with_jinja_hack_markers}"

          j2 "''${to_be_jinjad}" > "$(basename "''${f}" .stl).gcode"
          rm "''${to_be_jinjad}"
        done
      '';
    in
    {
      packages.x86_64-linux = {
        inherit convert-to-protoprinter;
        default = convert-to-protoprinter;
      };
      devShell.x86_64-linux = pkgs.mkShell {
        packages = [ convert-to-protoprinter ] ++ used-packages;
        shellHook = ''
          echo Usage: convert-to-protoprinter [stl file]...
        '';
      };
    };
}
