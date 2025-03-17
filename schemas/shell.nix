{ 
  pkgs,
  ...
}:

let
  schemasDependencies = with pkgs; [
    redocly
    just
    quicktype
    uv
  ];
in
pkgs.mkShell {
  buildInputs = schemasDependencies;
}
