{
  description = "A simple Go development shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            gcc
            go
            gopls
          ];
        };

        packages = {
          gamemode-auto-enable = pkgs.buildGoModule {
            pname = "gamemode-auto-enable";
            version = "0.0.1";
            src = self;
            vendorHash = "sha256-r9FG8gtQzkf/Idb1b8TrSoyoxBtwFTKosqeeaBRhiTQ=";
          };
          default = self.packages.${system}.gamemode-auto-enable;
        };

        nixosModules = {};
       });
}