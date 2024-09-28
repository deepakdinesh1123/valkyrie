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
  vendorHash = "sha256-WO08d65qux1CD9r6+E8uduSz8gMDDgKMo1di+E0sK4k=";
  
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