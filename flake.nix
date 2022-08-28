{
  description = "Prometheus labels injector";

  inputs = {
    hotPot.url = "github:shopstic/nix-hot-pot";
    nixpkgs.follows = "hotPot/nixpkgs";
    flakeUtils.follows = "hotPot/flakeUtils";
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.utils.follows = "flakeUtils";
    };
  };

  outputs = { self, nixpkgs, flakeUtils, hotPot, gomod2nix }:
    flakeUtils.lib.eachSystem [ "aarch64-darwin" "aarch64-linux" "x86_64-darwin" "x86_64-linux" ]
      (system:
        let
          pkgs = import nixpkgs {
            inherit system;
            overlays = [
              gomod2nix.overlays.default
            ];
          };
          hotPotPkgs = hotPot.packages.${system};
          vscodeSettings = pkgs.writeTextFile {
            name = "vscode-settings.json";
            text = builtins.toJSON {
              "yaml.schemaStore.enable" = true;
              "yaml.schemas" = {
                "https://json.schemastore.org/github-workflow.json" = ".github/workflows/*.yaml";
              };
              "nix.enableLanguageServer" = true;
              "nix.formatterPath" = pkgs.nixpkgs-fmt + "/bin/nixpkgs-fmt";
              "nix.serverPath" = pkgs.rnix-lsp + "/bin/rnix-lsp";
            };
          };
          promLabelsInjector = pkgs.callPackage ./build.nix { };
        in
        rec {
          defaultPackage = promLabelsInjector;
          packages = pkgs.lib.optionalAttrs pkgs.stdenv.isLinux {
            image = pkgs.callPackage ./image {
              inherit promLabelsInjector;
            };
          };
          devShell = pkgs.mkShellNoCC {
            shellHook = ''
              mkdir -p ./.vscode
              cat ${vscodeSettings} > ./.vscode/settings.json
            '';
            buildInputs = [ gomod2nix.packages.${system}.default ] ++ builtins.attrValues {
              inherit (hotPotPkgs)
                manifest-tool
                ;
              inherit (pkgs)
                skopeo
                awscli2
                go
                binutils
                ;
            };
          };
        }
      );
}
