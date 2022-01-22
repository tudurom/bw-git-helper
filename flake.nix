{
  description = "Bitwarden Git credential helper";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-21.11";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs@{ self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        defaultPackage = import ./default.nix { inherit pkgs; };
        defaultShell = import ./shell.nix { inherit pkgs; };
      }
    );
}
