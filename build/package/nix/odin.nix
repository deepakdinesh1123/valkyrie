{ lib,
  pkgs
}:


pkgs.buildGoModule rec {
  pname = "odin";
  version = "0.0.1";

  vendorHash = "sha256-DYU2DMM5OVUvrWx5AayhhQIqn6Pd/okjspnsiwB1c/w=";
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


