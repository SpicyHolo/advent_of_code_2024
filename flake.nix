{
  description = "A minimal Nix Flake for a Go development shell";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs";

  outputs = { self, nixpkgs }: {
    packages.x86_64-linux.default = let
    pkgs = import nixpkgs {
      system = "x86_64-linux";
    };

    in pkgs.mkShell {
      buildInputs = [
        pkgs.go
      ];

      shellHook = ''
        echo "Python environment with pandas and numpy activated!"
        export GOPATH=$PWD/.gopath
        export PATH=$GOPATH/bin:$PATH
      '';
    };
  };
}

