{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils = { url = "github:numtide/flake-utils"; };
    licenser = { url = "github:liamawhite/licenser/62520dbef14ff6e9aa864e0dbc19da9e3bed61c0"; };
  };
  outputs = { self, nixpkgs, utils, licenser }:
    utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        with pkgs;
        {
          devShells.default = pkgs.mkShell {
            buildInputs = [
              git
              go
              ffmpeg
              ttyd
              vhs
              licenser.packages.${system}.default
            ];
          };
        }
      );
}

