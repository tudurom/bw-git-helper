{
  description = "Bitwarden Git credential helper";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-21.11";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in rec {
        packages.bw-git-helper = let
          spdx = lic: lic // {
            url = "https://spdx.org/licenses/${lic.spdxId}.html";
          };
        in pkgs.buildGoModule rec {
          pname = "bw-git-helper";
          version = "unstable";

          src = ./.;

          vendorSha256 = "1x7iwd4ndcvcwyx8cxhkcwn5kwwf8nam47ln743wnvcrxi7q2604";
          buildInputs = with pkgs; [ bitwarden-cli ];

          meta = with pkgs.lib; {
            description = "A git credential helper using BitWarden as a backend";
            homepage = "https://github.com/tudurom/bw-git-helper";
            license = spdx { spdxId = "EUPL-1.2"; };
            maintainers = [ maintainers.tudor ];
            platforms = platforms.all;
          };
        };
        defaultPackage = packages.bw-git-helper;
        apps.bw-git-helper = utils.lib.mkApp { drv = packages.bw-git-helper; };
        defaultApp = apps.bw-git-helper;
    });
}
