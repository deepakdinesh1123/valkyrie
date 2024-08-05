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
  vendorHash = "sha256-g+YA2d4tuAtGazjtNiIyyaWbJfnZXMeHk7e8EDr+uUw=";
  
  src = ../../..;
  
  doCheck = false;
  
  buildInputs = [
    gpgme
    libgpg-error
    libassuan
  ] ++ lib.optionals stdenv.isLinux [
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