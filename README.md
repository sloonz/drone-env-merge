An environment variable plugin to attach mutliple environment variable plugins to a drone runner.

Note that all plugins uses the same secret.

_Please note this project requires Drone server version 1.6 or higher._

## Usage

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:80 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --env=DRONE_UPSTREAMS=http://drone-env-plugin-1,http://drone-env-plugin-2,http://drone-env-plugin-3 \
  --restart=always \
  --name=drone-env-merge
```

Update your runner configuration to include the plugin address and the shared secret.

```text
DRONE_ENV_PLUGIN_ENDPOINT=http://drone-env-merge:3000
DRONE_ENV_PLUGIN_TOKEN=bea26a2221fd8090ea38720fc445eca6
```
