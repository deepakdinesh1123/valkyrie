{ 
  pkgs,
  ...
}:

let
  docsDependencies = with pkgs; [
    python313
    uv
  ];
in
pkgs.mkShell {
  buildInputs = docsDependencies;
}
