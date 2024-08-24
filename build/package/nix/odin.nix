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
  vendorHash = "sha256-dofEXl3JiXq0Xqk6bGqW5gSpPfBo2Y5SH7g47Rt9ZDo=";
  
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
    license = licenses.mit;
    maintainers = with maintainers; [ deepak ];
    mainProgram = "odin";
  };
}