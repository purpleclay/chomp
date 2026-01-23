{
  description = "A parser combinator library for chomping strings (a rune at a time) in Go";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    go-overlay.url = "github:purpleclay/go-overlay";
  };

  outputs = { self, nixpkgs, flake-utils, go-overlay }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ go-overlay.overlays.default ];
        };
      in
      with pkgs;
      {
        devShells.default = mkShell {
          buildInputs = [
            git
            (go-bin.fromGoMod "${self}/go.mod")
            gofumpt
            golangci-lint
            go-task
          ];
        };
      }
    );
}
