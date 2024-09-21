{ lib
, pkgs
, stdenv
, buildGoModule
, pkg-config
, gpgme
, btrfs-progs
, libgpg-error
, libassuan
}:

buildGoModule rec {
  pname = "odin";
  version = "0.0.1";
  vendorHash = "sha256-BCy81Enful4Ajdo8krQKfb7DVGQssXK6FDld0URsm+c=";
  
  src = ../../..;
  
  doCheck = false;
  
  buildInputs =  lib.optionals stdenv.isLinux [
    gpgme
    libgpg-error
    libassuan
    btrfs-progs
  ];
  
  nativeBuildInputs = [ pkg-config ];
  
  subPackages = [ "cmd/odin" ];
  
  ldflags = [ "-s" "-w" "-X info.version=${version}" ];
  
  meta = with lib; {
    description = "Odin Server";
    license = licenses.asl20;
    maintainers = with maintainers; [ ];
    mainProgram = "odin";
  };
}