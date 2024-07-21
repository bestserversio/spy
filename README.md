A Go program that retrieves and updates game server information such as user counts, map names, and more from the Best Servers API along with adds/updates game servers using APIs such as the [Valve Master Server](https://developer.valvesoftware.com/wiki/Master_Server_Query_Protocol) and more!

## Building
You may use `git`, `make`, and the [`build.sh`](./build.sh) script to build/install this project. The build script results in the `spy` executable being outputted to the [`build/`](./build/) directory.

```bash
# Clone repository and change to its directory.
git clone https://github.com/bestserversio/spy

cd spy/

# Build Spy project.
./build.sh

# Install Spy executable to /usr/bin and configuration file to /etc/bestservers/spy.json.
sudo make install

# Run project.
spy
```

## Command Line
The following command line arguments are supported.

* `--list -l` => Prints the current configuration and exits.
* `--version -v` => Prints the current version and exits.
* `--help -h` => Prints the current supported command line arguments and exits.
* `--cfg` => The full path to the config file (default `/etc/bestservers/spy.json`).

<details>
    <summary>Example(s)</summary>

Load Spy using the local `local.json` configuration file.

```bash
spy --cfg=./local.json
```

Print help menu.
```bash
spy -h
```

Print current configuration.
```bash
spy -l
```
</details>

## Configuration
We use JSON to parse configuration files and the default configuration file is located at `/etc/bestservers/spy.json`. You may change this path with the `--cfg` command line argument described above.

### Types
The following describes the types used in the configuration documentation below. Please not types that are not used in the configuration are not listed.

| Name | Size (Bytes) | Description |
| ---- | ------------ | ----------- |
| string | N/A | Characters inside of quotes (`""`)
| int | 4 | A signed 32-bit integer. |
| uint | 4 | An unsigned 32-bit integer. |
| bool | 1 | A simple `true` or `false` type. |

### General
Please take a look at the following configuration.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| verbose | int | `0` | The verbose level. Values 1 - 6 are supported. |
| api | BsApi Object | `{}` | The Best Servers API object (read below). |
| web_api | Web API Object | `{}` | The web API object (read below). |
| vms | Valve Master Server Object | `{}` | The Valve Master Server object (read below). |
| scanners | Scanner Array | `[]` | An array of Scanner objects (read below). |
| platform_maps | Platform Maps Array | `[]` | An array of platform map objects (read below). |
| bad_names | []string | `[]` | An array of strings that represent bad names to be filtered. |
| bad_ips | []string | `[]` | An array of strings that represent bad IP ranges to be filtered. CIDR ranges are supported (`/24`). Examples include `192.168.3.0/24` and `192.168.3.5`. |
| bad_asns | []uint | `[]` | An array of unsigned integers that represent bad ASNs to be filtered (e.g. `AS<uint>`). |

<details>
    <summary>Example(s)</summary>

Using a verbose level of `5`.

```json
{
    "verbose": 5,
    "api": {

    },
    "web_api": {

    },
    "vms": {

    },
    "scanners": [

    ],
    "platform_maps": [

    ]
}
```

Setting a few bad names, IPs, and ASNs.

```json
{
    "bad_words": ["badword1", "badword2", "badword3"],
    "bad_ips": ["10.4.0.0/24", "1.2.3.4", "10.5.0.0/16"],
    "bad_asns": [12345, 43212, 42121]
}
```
</details>

### Bs API Object
The Best Servers API object includes connection and authorization details to the Best Servers API used for adding/updating servers.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| host | string | `"http://localhost"` | The host and endpoint for the Best Servers API. |
| authorization | string | `""` | The authorization token/string to set in the `Authorization` request header. |
| timeout | int | `5` | The request timeout. |

<details>
    <summary>Example(s)</summary>

Use API at `http://myserver.com/api/server/add` with the authorization `Bearer test` and a timeout of `10` seconds.

```json
{
    "host": "http://myserver.com/api/server/add",
    "authorization": "Bearer test",
    "timeout": 10
}
```
</details>

### Web API Project
The web API/polling information to retrieve configuration settings for this Spy instance from.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| enabled | bool | `false` | Whether to enable web API polling. |
| host | string | `"http://localhost"` | The web API host. |
| endpoint | string | `"/api/spy/get"` | The web API endpoint (comes after host in request). |
| authorization | string | `""` | What to set the `Authorization` request header to. |
| timeout | int | `5` | The API request timeout. |
| interval | int | `120` | How often in seconds to polling from the web config (0 = only poll once on Spy startup). |

<details>
    <summary>Example(s)</summary>

Use web API at `http://myserver.com/api/spy/get` with the authorization `Bearer test` and a timeout of `10` seconds. We want Spy to update every 1200 seconds with information from the web API.

```json
{
    "enabled": true,
    "host": "http://myserver.com",
    "endpoint": "/api/server/add",
    "authorization": "Bearer test",
    "timeout": 10,
    "interval": 1200
}
```
</details>

## VMS Object
The Valve Master Server API information.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| enabled | bool | `false` | Whether to enable VMS. |
| timeout | int | `5` | The request timeout. |
| api_token | string | `""` | The token/key to use when retrieving game servers from this API. |
| app_ids | []uint | `[]` | A list of Steam app IDs to retrieve for. Make sure the `platform_maps` array object contains a mapping from the app IDs to the platform IDs! |
| recv_only | bool | `false` | If enabled, Spy will retrieve game servers from VMS, but not update/add servers to the Best Servers API. |
| min_wait | int | `60` | The minimum amount of time in seconds to wait between each VMS request. |
| max_wait | int | `180` | The maximum amount of time in seconds to wait between each VMS request. |
| limit | int | `100` | The limit of game servers to retrieve per VMS request. |
| exclude_empty | bool | `true` | If true, will exclude empty servers from the VMS request directly. |
| sub_bots | bool | `true` | If true, subtracts the bot count from the user count when calcaluting the current user count. |
| add_only | bool | `false` | If true, game servers from each VMS request will be added to Best Servers only if it doesn't already exist. |

<details>
    <summary>Example(s)</summary>

Enable the Valve Master Server for app IDs `240` (Counter-Strike: Source) and `440` (Team Fortress 2). We want to request the VMS every `1000` - `2000` seconds to avoid rate-limiting. We also want to only add new servers to Best Servers and exclude empty servers from the VMS request.

```json
{
    "enabled": true,
    "api_token": "MY_STEAM_API_TOKEN",
    "min_wait": 1000,
    "max_wait": 2000,
    "exclude_empty": true,
    "add_only": true
}
```
</details>

### Scanner Object
This is the scanner object. Scanners query existing servers and update their information with different types of query types.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| protocol | string | `"A2S"` | The query protocol to use (supported types are `"A2S"` so far). |
| platform_ids | []uint | `[]` | An array of platform IDs to query with this scanner. |
| limit | int | `100` | The maximum amount of servers allowed to be query with each batch. |
| min_wait | int | `30` | The minimum amount of time in seconds to wait between each query batch. |
| max_wait | int | `60` | The maximum amount of time in seconds to wait between each query batch. |
| recv_only | bool | `false` | If true, servers will be scanned/outputted (depending on verbose level), but won' be updated through the Best Servers API. |
| sub_bots | bool | `false` | If true, will subtract the bot count from the user count when calculating the user count. |
| query_timeout | int | `3` | The query timeout in seconds. |
| a2s_player | bool | `true` | If the protocol is `A2S`, will attempt to send an `A2S_PLAYER` request alongside `A2S_INFO` to determine if the game server is online/not spoofed. |


<details>
    <summary>Example(s)</summary>

We want to create a scanner that queries existing servers with the `A2S` protocol and executes every `5` - `10` seconds. The limit will be `1000` servers per execution and we want to also send an `A2S_PLAYER` request to help prevent spoofed servers. We will be scanning platform IDs `1` (Counter-Strike: Source) and `2` (Team Fortress 2).

```json
[
    {
        "protocol": "A2S",
        "platform_ids": [1, 2],
        "limit": 1000,
        "min_wait": 5,
        "max_wait": 10,
        "a2s_player": true
    }
]
```
</details>

### Platform Maps Object
The platform maps object is used to map app Ids (from something like the VMS) to platform Ids.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| app_id | int | `NULL` | The app ID to from. |
| platform_id | int | `NULL` | The platform ID to map to. |

<details>
    <summary>Example(s)</summary>

We want to map app ID `240` to platform ID `1` (Counter-Strike: Source) and `440` to platform ID `2` (Team Fortress 2)

```json
[
    {
        "app_id": 240,
        "platform_id": 1
    },
    {
        "app_id": 440,
        "platform_id": 2
    }
]
```
</details>

## Credits
* [Christian Deacon](https://github.com/gamemann)