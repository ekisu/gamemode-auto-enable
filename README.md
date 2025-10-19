# gamemode-auto-enable

A small daemon that automatically enables GameMode when a game is running. Currently supports games launched through Steam.

# Installation

## On flake-enabled NixOS systems

Add this flake to your system's `flake.nix` inputs:

```nix
{
  inputs = {
    # Other inputs...
    gamemode-auto-enable.url = "github:ekisu/gamemode-auto-enable";
  };
}
```

Then, in your `nixosConfiguration`, add the module to your `modules` list and enable the service:

```nix
{
  # ...
  modules = [
    # Other modules...
    gamemode-auto-enable.nixosModules.default,
    {
      services.gamemode-auto-enable.enable = true;
    }
  ];
  # ...
}
```

## On other systems

First, clone the repository and build the binary:

```sh
git clone https://github.com/ekisu/gamemode-auto-enable.git
cd gamemode-auto-enable
make
```

Then, install the binary and the SystemD service:

```sh
sudo make install
```

Finally, enable the user service:

```sh
systemctl --user enable --now gamemode-auto-enable.service
```
