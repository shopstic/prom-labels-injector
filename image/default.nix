{ dumb-init
, dockerTools
, promLabelsInjector
, bash
}:
dockerTools.buildLayeredImage
{
  name = "prom-labels-injector";
  tag = promLabelsInjector.version;
  contents = [ bash dumb-init promLabelsInjector ];
  config = {
    Env = [
      "PATH=/bin"
    ];
    Entrypoint = [ "dumb-init" "--" "prom-labels-injector" ];
  };
}
