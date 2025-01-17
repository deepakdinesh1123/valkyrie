{ 
  pkgs,
  ...
}:

let
  schemasDependencies = with pkgs; [
    redocly
    just
  ];
in
pkgs.mkShell {
  buildInputs = schemasDependencies;
}
