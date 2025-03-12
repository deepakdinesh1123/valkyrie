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
  pname = "odin";
  version = "0.0.1";
  vendorHash = "sha256-jsrYRqnYSYipZ5KmT1LLwQbVZUcHGTwOKz5ww1jUM1Y=";

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

  subPackages = [ "cmd/odin" ];

  ldflags = [ "-s" "-w" "-X info.version=${version}" ];

  meta = with lib; {
    description = "Odin";
    license = licenses.asl20;
    maintainers = with maintainers; [ deepak sujay manoj ];
    mainProgram = pname;
  };
}
