import http from 'k6/http';
import { sleep } from 'k6';

var tasks = [
  {
    "file": {
      "name": "main.py",
      "content": "print('hello world')"
    },
    "environment": `
    {
      description = "A simple flake example";

      inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
      inputs.flake-utils.url = "github:numtide/flake-utils";

      outputs = { self, nixpkgs, flake-utils, ... }:
        flake-utils.lib.eachDefaultSystem (system: 
            let
                pkgs = import nixpkgs { inherit system; };
            in
            rec
            {
                packages = {
                    something = pkgs.writeShellApplication {
                        name = "something";
                        runtimeInputs = [ pkgs.python3 ];
                        text = ''
                            python script.py
                        '';
                    };
                };
                apps.default = {
                    type = "app";
                    program = "\${packages.something}/bin/something";
                };
            }
        );
      }
    `
  },
  {
    "file": {
      "name": "main.py",
      "content": "print('hello world')"
    },
    "environment": `
    {
      description = "A simple flake example";

      inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
      inputs.flake-utils.url = "github:numtide/flake-utils";

      outputs = { self, nixpkgs, flake-utils, ... }:
        flake-utils.lib.eachDefaultSystem (system: 
            let
                pkgs = import nixpkgs { inherit system; };
            in
            rec
            {
                packages = {
                    something = pkgs.writeShellApplication {
                        name = "something";
                        runtimeInputs = [ pkgs.python39 ];
                        text = ''
                            python script.py
                        '';
                    };
                };
                apps.default = {
                    type = "app";
                    program = "\${packages.something}/bin/something";
                };
            }
        );
      }
    `
  },
];

export const options = {
  // A number specifying the number of VUs to run concurrently.
  vus: 10,
  // A string specifying the total duration of the test run.
  duration: '300s',
};

export default function() {
  http.get('https://test.k6.io');
  sleep(1);
}
