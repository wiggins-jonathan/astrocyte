{
  description = "Astrocyte development environment";

  nixConfig.bash-prompt-prefix = "(✶ ) ";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.05";
  };

  outputs = { self, nixpkgs }:
    let
      systems = ["x86_64-linux" "x86_64-darwin"];
      forEachSystem = f: nixpkgs.lib.genAttrs systems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in {
      devShells = forEachSystem ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
          ];

          shellHook = ''
            echo "$(go version)"
          '';
        };
      });
    };
}
