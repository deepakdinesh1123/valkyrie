{ lib,
  pkgs
}:


pkgs.buildGoModule rec {
  pname = "odin";
  version = "0.0.1";

  vendorHash = "sha256-pu9Pz7wmmCzxryg2bXl9qSeH22gLBrKR1v0W9V85naY=";
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


