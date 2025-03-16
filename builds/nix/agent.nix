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
  vendorHash = "sha256-5B9kZUCDIM9hXRaMZcVUDHa7H4Qn7D9iUl+J51LPoOE=";

  src = ../../agent;

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
