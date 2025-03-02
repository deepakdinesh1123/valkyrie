{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    gcc libunistring tzdata mailcap iana-etc libidn2 glibc go
  ];
}
