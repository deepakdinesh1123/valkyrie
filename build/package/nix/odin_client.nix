{ lib
, pkgs
, stdenv
, buildGoModule
}:

buildGoModule rec {
  pname = "odinc";
  version = "0.0.1";
  vendorHash = "sha256-dofEXl3JiXq0Xqk6bGqW5gSpPfBo2Y5SH7g47Rt9ZDo=";
  
  src = ../../..;
  
  doCheck = false;
  
  subPackages = [ "pkg/odin/odinc" ];
  
  ldflags = [ "-s" "-w" "-X info.version=${version}" ];
  
  meta = with lib; {
    description = "Odin Client";
    license = licenses.mit;
    maintainers = with maintainers; [ deepak ];
    mainProgram = "odinc";
  };
}