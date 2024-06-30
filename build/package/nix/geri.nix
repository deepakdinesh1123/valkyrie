{ lib,
  pkgs
}:


pkgs.buildGoModule rec {
  pname = "geri";
  version = "0.0.1";

  vendorHash = "sha256-KHmpKCc3OzTZKDOGkYcud9FdD7mLjbfcu82ReboUtCE=";
  doCheck = false;

  src = ../../..;

  ldflags = [ "-s" "-w" "-X info.version=${version}" ];

  subPackages = [ "cmd/geri" ];

  meta = with lib; {
    description = "geri";
    license = licenses.asl20;
    maintainers = with maintainers; [ ];
    mainProgram = "geri";
  };
}


