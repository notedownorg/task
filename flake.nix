{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    licenser = { url = "github:liamawhite/licenser/bdf2c1beeaf09aae9f27b9d680b5b6aa22e4f1a0"; };
  };
  outputs = { self, nixpkgs, utils, licenser }:
    utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          devShells.default = pkgs.mkShell {
            buildInputs = [
              pkgs.git
              pkgs.go
              pkgs.ffmpeg
              pkgs.ttyd
              pkgs.vhs
              licenser.packages.${system}.default
            ];
          };
        }
      );
}

