{ pkgs }:
with pkgs;

# https://nixos.org/manual/nixpkgs/stable/#sec-language-go
buildGoModule rec {
  pname = "licenser";
  version = "v0.7.0";

  src = fetchFromGitHub {
    owner = "liamawhite";
    repo = "licenser";
    rev = "${version}";
    sha256 = "sha256-NdK2wM7ivGV/tenJuOPuzinnKMiNTOJe0QhRMjzM3JU=";
  };

  vendorHash = "sha256-LBVVhg69VdQVsVARCkwooe6N6DacgViIW/iQWHCya4k=";
  ldFlags = "-w -s";
}

