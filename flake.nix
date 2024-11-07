{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        licenser = import ./nix/packages/licenser.nix { inherit pkgs; };
      in
      with pkgs;
      {
        devShells.default = mkShell {
          buildInputs = [
            git
            go
            ffmpeg
            ttyd
            vhs
            licenser
          ];
        };
      }
    );
}

