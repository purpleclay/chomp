{
  lib,
  buildGoApplication,
  go-bin,
}: let
  examples = [
    "git-diff"
    "gpg"
  ];

  mkExample = name:
    buildGoApplication {
      pname = "chomp-example-${name}";
      version = "0.1.0";
      src = ./examples;
      go = go-bin.fromGoMod ./examples/go.mod;
      modules = ./examples/govendor.toml;
      subPackages = [name];

      # Provide the chomp library source for the local replace directive
      # (go.mod has: replace github.com/purpleclay/chomp => ../)
      localReplaces = {
        "github.com/purpleclay/chomp" = ./.;
      };

      meta = with lib; {
        description = "Chomp parser combinator example: ${name}";
        homepage = "https://github.com/purpleclay/chomp";
        license = licenses.mit;
      };
    };
in
  lib.genAttrs examples mkExample
