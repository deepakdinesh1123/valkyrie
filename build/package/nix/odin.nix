{ lib,
  pkgs,
  stdenv
}:


pkgs.buildGoModule rec {
  pname = "odin";
  version = "0.0.1";

  vendorHash = "sha256-g+YA2d4tuAtGazjtNiIyyaWbJfnZXMeHk7e8EDr+uUw=";
  doCheck = false;

  buildInputs = lib.optionals stdenv.isLinux [ pkgs.gpgme ];
  nativebuildInputs = [ pkgs.pkg-config ];

  src = ../../..;

  subPackages = [ "cmd/odin" ];
  ldflags = [ "-s" "-w" "-X info.version=${version}" ];

  meta = with lib; {
    description = "Odin Server";
    license = licenses.asl20;
    maintainers = with maintainers; [ ];
    mainProgram = "odin";
  };
}


