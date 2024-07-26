{ lib,
  pkgs
}:


pkgs.buildGoModule rec {
  pname = "odin";
  version = "0.0.1";

  vendorHash = "sha256-tCmvn1eyblZSFoDBJ0xTSyBaZEbe7FtsXjtTNlZkKoc=";
  doCheck = false;

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


