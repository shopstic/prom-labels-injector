{ lib
, fetchurl
, runCommandNoCC
, buildGoApplication
}:
buildGoApplication rec {
  pname = "prom-labels-injector";
  version = "1.0.0";
  src = ./src;
  modules = ./src/gomod2nix.toml;
}

