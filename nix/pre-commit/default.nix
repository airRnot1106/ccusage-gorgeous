{ config }:
{
  check.enable = true;
  settings.hooks = {
    golangci-lint.enable = true;
    gotest.enable = true;
    nil.enable = true;
    treefmt = {
      enable = true;
      package = config.treefmt.build.wrapper;
    };
  };
}
