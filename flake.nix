{
  description = "University week counter API.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    inputs:
    with inputs;
    {
      nixosModules.default = import ./module.nix self.outputs.packages;
    }
    //
      flake-utils.lib.eachSystem
        [
          "x86_64-linux"
          "aarch64-linux"
        ]
        (
          system:
          let
            pkgs = import nixpkgs { inherit system; };
            version = builtins.substring 0 8 self.lastModifiedDate or "dirty";
            uni-week-counter = pkgs.callPackage ./package.nix { inherit version; };
          in
          {
            packages.uni-week-counter = uni-week-counter;
            packages.default = uni-week-counter;
            devShells.default = pkgs.mkShell { packages = with pkgs; [ go ]; };
            formatter = pkgs.nixfmt-rfc-style;
          }
        );
}

