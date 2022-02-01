{ dumb-init
, dockerTools
, promLabelsInjector
}:
dockerTools.buildLayeredImage
{
  name = "prom-labels-injector";
  config = {
    Entrypoint = [ "${dumb-init}/bin/dumb-init" "--" "${promLabelsInjector}/bin/prom-labels-injector" ];
  };
}
