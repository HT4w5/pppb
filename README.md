# pppcm

Concurrently start ppp links and periodically check link state. Restart links if down.

## Prerequisites
- Linux
- [`pppd`](https://github.com/ppp-project/ppp)

> Tested with `Rocky Linux 9.6 x86_64` and `pppd version 2.4.9`.

## Features
- Start all links concurrently.
- Check number of running links at regular intervals.
- Restart all links if less than a certain number of links are alive.
- Force restart all links at regular intervals.

## Usage

### `pppcm`

- Reads `config.json` from CWD.
- Config format:
    ```json
    {
        "links": [
            {
                "tag": "isp0",
                "ttyname": "eth0",
                "user": "<user>",
                "password": "<password>",
                "ifname": "ppp0"
            },
            {
                "tag": "isp1",
                "ttyname": "eth1",
                "user": "<user>",
                "password": "<password>",
                "ifname": "ppp1"
            },
            {
                "tag": "isp2",
                "ttyname": "eth2",
                "user": "<user>",
                "password": "<password>",
                "ifname": "ppp2"
            }
        ],
        "daemon": {
            "run_dir": "/var/run",
            "expected": 3,
            "enabled": true,
            "check_interval": 300,
            "force_restart": false,
            "force_restart_interval": 86400
        }
    }
    ```
- `links`
    - `tag`      Unique tag of link.
    - `ttyname`  Device used by PPP link.
    - `user`     Name for authentication.
    - `password` Password for authentication.
    - `ifname`   Name of PPP network interface.
- `daemon`
    - `enabled` Enable daemon mode.
    - `run_dir` Runtime variable directory - where `[ifname].pid` files are located.
    - `expected` Expected number of links. Will trigger restart if less are alive.
    - `check_interval` Check interval in seconds. Minimum `300`.
    - `force_restart` Enable force restart.
    - `force_restart_interval` Force restart interval in seconds. Must be greater than `check_interval`.

### `pppcm-cfg-gen`

- Generates JSON config for many link with the same account.
- Prints to `stdout`.
    ```console
    Usage: pppcm-cfg-gen <tag-prefix> <ttyname-prefix> <user> <password> <ifname-prefix> <link-count>
    ```

## License
[AGPL-3.0](LICENSE)