{ lib
, pkgs
, stdenv
, buildGoModule
, pkg-config
, gpgme
, btrfs-progs
, libgpg-error
, libassuan
}:

buildGoModule rec {
  pname = "valkyrie";
  version = "0.0.1";
  vendorHash = "sha256-bqWapZcdV/2jyzB8PCesOtrZ6IacyyMmx5GK52sCEuE=";

  src = ../..;

  doCheck = false;

  buildInputs = lib.optionals stdenv.isLinux [
    gpgme
    libgpg-error
    libassuan
    btrfs-progs
  ];

  nativeBuildInputs = [ pkg-config ];

  tags = lib.optionals stdenv.isDarwin [ "darwin" ]
      ++ lib.optionals stdenv.isLinux [ "docker" ];

  subPackages = [ "cmd/valkyrie" ];

  ldflags = [ "-s" "-w" "-X info.version=${version}" ];

  meta = with lib; {
    description = "valkyrie";
    license = licenses.asl20;
    maintainers = with maintainers; [ deepak sujay manoj ];
    mainProgram = pname;
  };
}
