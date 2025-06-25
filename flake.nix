{
  inputs = {
    devenv.url = "github:cachix/devenv";
    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    nixpkgs.url = "github:cachix/devenv-nixpkgs/rolling";
    systems.url = "github:nix-systems/default";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    { flake-parts, ... }@inputs:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.git-hooks.flakeModule
        inputs.devenv.flakeModule
        inputs.treefmt-nix.flakeModule
      ];
      systems = import inputs.systems;

      perSystem =
        {
          config,
          lib,
          pkgs,
          system,
          ...
        }:
        {
          packages.default = pkgs.buildGoModule {
            pname = "ccugorg";
            version = "0.0.1";
            src = lib.cleanSource ./.;
            vendorHash = "sha256-ZYY5VCbGmSQWhcAaPXMFgq4LWOa9jQ/gDt4zJjLzBco=";

            postInstall = ''
              mv $out/bin/ccusage-gorgeous $out/bin/ccugorg
            '';
          };

          devenv.shells.default = {
            packages = with pkgs; [
              git
              nil
            ];
            containers = lib.mkForce { };
            languages.go.enable = true;
            enterShell = ''
              ${config.pre-commit.installationScript}
            '';
          };

          pre-commit = import ./nix/pre-commit {
            inherit config;
            inherit pkgs;
          };
          treefmt = import ./nix/treefmt;
        };
    };
}
