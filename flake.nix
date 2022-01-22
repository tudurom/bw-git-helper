{
  description = "Bitwarden Git credential helper";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-21.11";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs@{ self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in rec {
        packages = utils.lib.flattenTree {
          bw-git-helper = import ./default.nix { inherit pkgs; };
        };
        defaultPackage = packages.bw-git-helper;

        apps.bw-git-helper = utils.lib.mkApp { drv = packages.bw-git-helper; };
        defaultApp = apps.bw-git-helper;

        devShell = import ./shell.nix { inherit pkgs; };
      }
    );
}
