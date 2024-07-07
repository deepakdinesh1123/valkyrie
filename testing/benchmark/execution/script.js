import http from 'k6/http';
import { sleep, check } from 'k6';

function getRandomItemFromArray(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}

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
                            python main.py
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
      "content": "import time\ntime.sleep(10)"
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
                            python main.py
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
  let res = http.post('http://localhost:8000/execution/execute/', JSON.stringify(getRandomItemFromArray(tasks)));
  const respData = res.json();
  check(res, {
    'status was 200': (r) => r.status == 200,
    'execution_id_exists': (r) => respData.execution_id != null,
  });
  sleep(1);
}
