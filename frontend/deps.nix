{ 
  pkgs,
  ...
}:
let
  frontendDependencies = with pkgs; [
    nodejs_22
  ];
in
 frontendDependencies
