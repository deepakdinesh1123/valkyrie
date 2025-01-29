{ lib
, pkgs
, stdenv
, buildGoModule
, pkg-config
, ...
}:

buildGoModule rec {
  pname = "agent";
  version = "0.0.1";
  vendorHash = "sha256-Nnyd+Jnc80FjNOlNV4oAY6vu9KrHkCs60I9+B8jKsMc=";

  src = ../../../agent;

  doCheck = false;

  nativeBuildInputs = [ pkg-config ];

  ldflags = [ "-s" "-w" "-X info.version=${version}" ];

  meta = with lib; {
    description = "Valkyrie Agent";
    license = licenses.asl20;
    maintainers = with maintainers; [ deepak sujay manoj ];
    mainProgram = pname;
  };
}
