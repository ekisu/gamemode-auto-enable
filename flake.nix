{
  description = "A simple Go development shell";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, flake-parts }@inputs:
    flake-parts.lib.mkFlake { inherit inputs; } (top@{ config, withSystem, moduleWithSystem, ... }: {
      flake = {
        nixosModules = {
          gamemode-auto-enable = moduleWithSystem (
            perSystem@{ pkgs, system, ... }:
            nixos@{ config, lib, ... }: {
              options.services.gamemode-auto-enable = {
                enable = lib.mkEnableOption "gamemode-auto-enable user service";
              };

              config = lib.mkIf config.services.gamemode-auto-enable.enable {
                systemd.user.services.gamemode-auto-enable = {
                  description = "Automatically enable GameMode when a game is running.";
                  partOf = [ "graphical-session.target" ];
                  wantedBy = [ "graphical-session.target" ];
                  serviceConfig = {
                    ExecStart = "${self.packages.${system}.gamemode-auto-enable}/bin/gamemode-auto-enable";
                    Restart = "on-failure";
                    RestartSec = "5s";
                  };
                };
              };
            }
          );

          default = self.nixosModules.gamemode-auto-enable;
        };
      };

      perSystem = { config, pkgs, self', ... }: {
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
          default = self'.packages.gamemode-auto-enable;
        };
      };

      systems = [ "x86_64-linux" ];
    });
        
}