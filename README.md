
# Docker machine driver for [ArvanCloud](https://www.arvancloud.com/)


## Usage

    $ docker-machine create \
      --driver arvan \
      --arvan-api-token=<YOU_API_TOKEN> \
      ar-ins-1

## Options

| Parameter                    | Env                    |
| ---------------------------- | ---------------------- |
| **`--arvan-api-token`**      | `ARVAN_API_TOKEN`      |
| **`--arvan-image`**          | `ARVAN_IMAGE`          |
| **`--arvan-region`**         | `ARVAN_REGION`         |
| **`--arvan-server-flavor`**  | `ARVAN_SERVER_FLAVOR`  |
| **`--arvan-network`**        | `ARVAN_NETWORK`        |
| **`--arvan-security-group`** | `ARVAN_SECURITY_GROUP` |
| **`--arvan-ssh-user`**       | `ARVAN_SSH_USER`       |

## License

MIT © [Amir Keshavarz](https://github.com/satrobit)