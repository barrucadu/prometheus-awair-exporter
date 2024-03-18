{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    gitignore = {
      url = "github:hercules-ci/gitignore.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, gitignore }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
      };
    in
    {
      formatter.${system} = pkgs.nixpkgs-fmt;

      devShells.${system}.default = pkgs.mkShell {
        packages = [ pkgs.go ];
      };

      packages.${system}.default = pkgs.buildGoModule rec {
        pname = "prometheus-awair-exporter";
        version = "0.0.0";

        src = gitignore.lib.gitignoreSource ./.;
        vendorHash = "sha256-XDHUotANzNtmm5C9A3Ccav90zi3AYYz7zxfGUEJskec=";

        doCheck = false;

        meta = {
          description = "A simple Awair exporter for Prometheus.";
          homepage = "https://github.com/barrucadu/prometheus-awair-exporter";
        };
      };
    };
}
