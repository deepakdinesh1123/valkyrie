{ 
  pkgs,
  ...
}:

let
  docsDependencies = with pkgs.python312Packages; [
    mkdocs-material
    mkdocs-minify-plugin
    markdown
  ];
in
pkgs.mkShell {
  buildInputs = docsDependencies;
}
