{ lib,
  pkgs
}:

let
  version = "1.68.0";
in
pkgs.buildGoModule rec {
  pname = "nardump";
  inherit version;

  vendorHash = "sha256-Hd77xy8stw0Y6sfk3/ItqRIbM/349M/4uf0iNy1xJGw=";
  doCheck = false;

  src = pkgs.fetchFromGitHub {
    owner = "tailscale";
    repo = "tailscale";
    rev = "v${version}";
    hash = "sha256-ETBca3qKO2iS30teIF5sr/oyJdRSKFqLFVO3+mmm7bo=";
  };

  subPackages = [ "cmd/nardump" ];

  ldflags = [
    "-w"
    "-s"
  ];
  meta = with lib; {
    homepage = "https://tailscale.com";
    description = "nardump";
    license = licenses.bsd3;
    mainProgram = "nardump";
    maintainers = with maintainers; [ mbaillie jk mfrw ];
  };
}


