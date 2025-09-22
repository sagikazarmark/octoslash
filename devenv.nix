{ pkgs, inputs, ... }:

{
  cachix.pull = [ "sagikazarmark-dev" ];

  overlays = [
    (final: prev: {
      dagger = inputs.dagger.packages.${final.system}.dagger;
    })
  ];

  languages = {
    go = {
      enable = true;
      package = pkgs.go_1_25;
    };
  };

  packages = with pkgs; [
    just
    golangci-lint
    goreleaser
    dagger
  ];
}
