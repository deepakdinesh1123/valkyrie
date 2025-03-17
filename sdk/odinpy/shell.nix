{ 
  pkgs, 
  ...
}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    uv
    just
  ];
}
