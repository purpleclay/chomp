{
  description = "A parser combinator library for chomping strings (a rune at a time) in Go";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    go-overlay = {
      url = "github:purpleclay/go-overlay";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };

    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    go-overlay,
    git-hooks,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [go-overlay.overlays.default];
        };

        pre-commit-check = git-hooks.lib.${system}.run {
          src = ./.;
          package = pkgs.prek;
          hooks = {
            alejandra = {
              enable = true;
              settings = {
                check = true;
              };
            };

            typos = {
              enable = true;
              entry = "${pkgs.typos}/bin/typos";
            };
          };
        };
      in
        with pkgs; {
          devShells.default = mkShell {
            inherit (pre-commit-check) shellHook;

            buildInputs =
              [
                alejandra
                (go-bin.fromGoMod "${self}/go.mod")
                gofumpt
                golangci-lint
                nil
                typos
              ]
              ++ pre-commit-check.enabledPackages;
          };
        }
    );
}
