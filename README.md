# Docker machine driver for [ArvanCloud](https://www.arvancloud.com)
> This library adds the support for creating [Docker machines](https://github.com/docker/machine) hosted on the [ArvanCloud](https://www.arvancloud.com).

[Click here](https://npanel.arvancloud.com/profile/api-keys) to create an API Token which needs to be put in the `--arvan-api-token` option.
  
## Installation

You can find sources and pre-compiled binaries [here](https://github.com/satrobit/docker-machine-driver-arvan/releases).

```bash
# Download the binary (this example downloads the binary for linux amd64)
$ wget https://github.com/satrobit/docker-machine-driver-arvan/releases/download/v0.1-alpha/docker-machine-driver-arvan_v0.1-alpha_linux_amd64.tar.gz
$ tar -xvf docker-machine-driver-arvan_v0.1-alpha_linux_amd64.tar.gz

# Make it executable and copy the binary in a directory accessible with your $PATH
$ chmod +x docker-machine-driver-arvan
$ cp docker-machine-driver-arvan /usr/local/bin/
```
## Usage

    $ docker-machine create \
      --driver arvan \
      --arvan-api-token=<YOU_API_TOKEN> \
      ar-ins-1

## Options

| Parameter                    | Env                    | Description |
| ---------------------------- | ---------------------- | ----------  |
| **`--arvan-api-token`**      | `ARVAN_API_TOKEN`      | API Token |
| **`--arvan-image`**          | `ARVAN_IMAGE`          | The linux image used to create an instance |
| **`--arvan-region`**         | `ARVAN_REGION`         | The region which the instance is located in |
| **`--arvan-server-flavor`**  | `ARVAN_SERVER_FLAVOR`  | The flavor (size) used for the instance
| **`--arvan-network`**        | `ARVAN_NETWORK`        | The network connected to the instances iface |
| **`--arvan-security-group`** | `ARVAN_SECURITY_GROUP` | Security group used for the instance |
| **`--arvan-ssh-user`**       | `ARVAN_SSH_USER`       | The SSH username that docker-machine tries to connect to |
| **`--arvan-ssh-port`**       | `ARVAN_SSH_PORT`       | The SSH port that docker-machine tries to connect to |

## License

MIT © [Amir Keshavarz](https://github.com/satrobit)
