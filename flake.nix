{
  description = "Bitwarden Git credential helper";

  inputs = {
    nixpkgs.url = "nixpkgs/release-20.09";
  };

  outputs = { self, nixpkgs }: let
    lib = nixpkgs.lib;

    systems = [
      "x86_64-linux"
      "i686-linux"
      "x86_64-darwin"
      "aarch64-linux"
      "armv6l-linux"
      "armv7l-linux"
    ];

    forAllSystems = f: lib.genAttrs systems (system: f system);
  in
    {
      defaultPackage = forAllSystems (system:
        let
          spdx = lic: lic // {
            url = "https://spdx.org/licenses/${lic.spdxId}.html";
          };

          pkgs = import nixpkgs {
            inherit system;
          };
        in
          pkgs.buildGoModule rec {
            pname = "bw-git-helper";
            version = "unstable";

            src = ./.;

            vendorSha256 = "1x7iwd4ndcvcwyx8cxhkcwn5kwwf8nam47ln743wnvcrxi7q2604";
            buildInputs = with pkgs; [ bitwarden-cli ];

            meta = with pkgs.stdenv.lib; {
              description = "A git credential helper using BitWarden as a backend";
              homepage = "https://github.com/tudurom/bw-git-helper";
              license = spdx { spdxId = "EUPL-1.2"; };
              maintainers = [ maintainers.tudor ];
              platforms = platforms.all;
            };
          }
      );
    };
}
