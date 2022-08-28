{ lib
, fetchurl
, runCommandNoCC
, buildGoApplication
}:
buildGoApplication rec {
  pname = "prom-labels-injector";
  version = "1.0.1";
  src = ./src;
  modules = ./src/gomod2nix.toml;
}

