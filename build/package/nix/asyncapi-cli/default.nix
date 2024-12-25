{ 
  lib
, stdenv
, fetchFromGitHub
, buildNpmPackage
, makeWrapper
, git
, chromium
, ...
}:

buildNpmPackage rec {
  pname = "cli";
  version = "2.14.0";

  src = fetchFromGitHub {
    owner = "asyncapi";
    repo = "cli";
    rev = "v${version}";
    hash = "sha256-HECJelflnfAJ0rTHsu+X5QgazxZdG8Ck2Jyv5ec2Q00=";
  };

  npmDepsHash = "sha256-VVemfHMsM1asxLM+s1MFEzCtinzrj0bVsRKDQaBcjT0=";

  npmBuildScript = "build";

  patches =  [./remove_example.diff];

  postPatch =  ''
    rm -rf src/commands/new
  '';



  # postInstall = ''
  #   # Remove the unnecessary binary symlink from the build
  #   rm -f $out/lib/node_modules/@asyncapi/bin/run

  #   cp -R packages/asyncapi $out/lib/node_modules/@asyncapi/

  #   # Create the binary wrapper
  #   mkdir -p $out/bin
  #   makeWrapper $out/lib/node_modules/@asyncapi/bin/run \
  #     $out/bin/asyncapi
  # '';

  # Specify native build inputs
  nativeBuildInputs = [ git chromium ];

  # Disable the default npm build step
  # dontNpmBuild = true;

  meta = with lib; {
    description = "CLI to work with your AsyncAPI files. You can validate them, use a generator, and bootstrap new files. Contributions are welcome.";
    homepage = "https://github.com/asyncapi/cli";
    changelog = "https://github.com/asyncapi/cli/blob/v${version}/CHANGELOG.md";
    license = licenses.asl20;
    maintainers = with maintainers; [ ];
    mainProgram = "asyncapi";
    platforms = platforms.all;
  };
}
