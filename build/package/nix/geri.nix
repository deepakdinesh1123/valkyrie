{ lib,
  pkgs
}:


pkgs.buildGoModule rec {
  pname = "geri";
  version = "0.0.1";

  vendorHash = "sha256-ssF8q2gS90lVo3zc42ShKQm0tz5J/nG2bS8ka30F8i8=";
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


