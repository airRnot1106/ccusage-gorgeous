{ config, pkgs }:
{
  check.enable = true;
  settings.hooks = {
    golangci-lint = {
      enable = true;
      extraPackages = with pkgs; [
        go
      ];
    };
    gotest.enable = true;
    nil.enable = true;
    treefmt = {
      enable = true;
      package = config.treefmt.build.wrapper;
    };
  };
}
