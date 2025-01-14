{ 
  pkgs,
  ...
}:

let
  oasDependencies = with pkgs; [
    redocly
  ];
in
pkgs.mkShell {
  buildInputs = oasDependencies;
}
