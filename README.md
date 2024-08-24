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
| log_directory | string | `./logs` | The path to the directory to store logs in. Use `NULL` to disable logging to files. |
| api | BsApi Object | `{}` | The Best Servers API object (read below). |
| web_api | Web API Object | `{}` | The web API object (read below). |
| vms | Valve Master Server Array | `[]` | An array of scanners for the Valve Master Server (read below). |
| scanners | Scanner Array | `[]` | An array of Scanner objects (read below). |
| platform_maps | Platform Maps Array | `[]` | An array of platform map objects (read below). |
| bad_names | []string | `[]` | An array of strings that represent bad names to be filtered. |
| bad_ips | []string | `[]` | An array of strings that represent bad IP ranges to be filtered. CIDR ranges are supported (`/24`). Examples include `192.168.3.0/24` and `192.168.3.5`. |
| bad_asns | []uint | `[]` | An array of unsigned integers that represent bad ASNs to be filtered (e.g. `AS<uint>`). |
| good_ips | []string | `[]` | A list of IP range(s) to avoid filtering with. |
| remove_inactive | Remove Inactive Object | `{}` | Remove inactive servers settings. |
| platform_filters | Platform Filter Array | `[]` | A list of platform-specific filters to apply. |
| remove_dups | Remove Duplicates Object | `{}` | A list of settings for removing duplicate servers by IP. |
| remove_timed_out | Remove Timed Out Object | `{}` | A list of settings for removing timed out servers. |

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
    "vms": [

    ],
    "scanners": [

    ],
    "platform_maps": [

    ],
    "bad_words": [

    ],
    "bad_ips": [

    ],
    "bad_asns": [

    ],
    "good_ips": [

    ],
    "platform_filters": [

    ],
    "remove_inactive": {

    },
    "remove_dups": {

    },
    "remove_timed_out": {

    }
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

### BS API Object
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
| save_to_fs | bool | `true` | Whether to save the web API to the local config file on the file system. |

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
    "interval": 1200,
    "save_to_fs": true
}
```
</details>

### VMS Object
The Valve Master Server API information.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| timeout | int | `5` | The request timeout. |
| api_token | string | `""` | The token/key to use when retrieving game servers from this API. |
| app_ids | []uint | `[]` | A list of Steam app IDs to retrieve for. Make sure the `platform_maps` array object contains a mapping from the app IDs to the platform IDs! |
| recv_only | bool | `false` | If enabled, Spy will retrieve game servers from VMS, but not update/add servers to the Best Servers API. |
| min_wait | int | `60` | The minimum amount of time in seconds to wait between each VMS request. |
| max_wait | int | `180` | The maximum amount of time in seconds to wait between each VMS request. |
| limit | int | `100` | The limit of game servers to retrieve per VMS request. |
| exclude_empty | bool | `true` | If true, will exclude empty servers from the VMS request directly. |
| only_empty | bool | `false` | If true, only empty servers will be returned by the VMS request. |
| sub_bots | bool | `true` | If true, subtracts the bot count from the user count when calcaluting the current user count. |
| add_only | bool | `false` | If true, game servers from each VMS request will be added to Best Servers only if it doesn't already exist. |
| random_apps | bool | `false` | If true, the next app ID is selected randomly instead of in order. |
| set_offline | bool | `true` | If true, all servers added by VMS will be set to offline by default. |
| update_limit | int | `0` | If above 0, will limit the amount of servers to update to this amount. |
| randomize_res | bool | `false` | If true, the results from the VMS request will be randomized. |

<details>
    <summary>Example(s)</summary>

Add a scanner for the Valve Master Server for app IDs `240` (Counter-Strike: Source) and `440` (Team Fortress 2). We want to request the VMS every `1000` - `2000` seconds to avoid rate-limiting. We also want to only add new servers to Best Servers and exclude empty servers from the VMS request.

```json
[
    {
        "api_token": "MY_STEAM_API_TOKEN",
        "min_wait": 1000,
        "max_wait": 2000,
        "exclude_empty": true,
        "add_only": true,
        "random_apps": false,
        "set_offline": true
    }
]
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
| random_platforms | bool | `false` | If true, the next platform ID will be selected randomly instead of in order. |
| visible_skip_count | int | `0` | If over 0, will only scan visible servers up to this value. When met, the counter is reset and the query scans for all servers including invisible. |
| request_delay | int | `0` | The delay in milliseconds between every server request. |

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

### Platform Map Object
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

### Remove Inactive Object
The remove inactive object contains settings on automatically removing inactive servers based off of last online.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| enabled | bool | `false` | Whether to enable automatically removing servers. |
| inactive_time | int | `2592000` | How many seconds since last online to consider a server inactive. |
| interval | int | `86400` | How often to run the remove inactive functionality in seconds. |
| timeout | int | `5` | The request timeout. |

<details>
    <summary>Example(s)</summary>

If we want to remove inactive servers that haven't been online in 500 seconds, we'd use the following. We'll execute the remove inactive functionality every 60 seconds (minute).

```json
{
    "enabled":  true,
    "inactive_time": 500,
    "interval": 60,
    "timeout": 5
}
```
</details>

### Platform Filter Object
The platform filters object is used to apply platform-specific filters to servers.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| id | int | `0` | The platform ID. |
| max_cur_users | int | `NULL` | The maximum current users allowed by servers in this platform. |
| max_users | int | `NULL` | The maximum max users allowed by servers in this platform. |
| allow_user_overflow | bool | `NULL` | If true, any servers that have more current users than max users will be filtered. |

<details>
    <summary>Example(s)</summary>

Let's say we want to limit Counter-Strike: Source servers to a maximum of 65 users (both current and max). The CS:S platform ID is `4`.

```json
[
    {
        "id": 4,
        "max_cur_users": 65,
        "max_users": 65
    }
]
```
</details>

### Remove Duplicates Object
This object contains settings for removing duplicate servers by IP. There are many fake game servers that use every port available on a single IP and these settings mitigate this issue.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| enabled | bool | `false` | Whether to enable remove duplicates checking. |
| interval | int | `120` | How often to check for duplicates in seconds. |
| limit | int | `100` | The amount of server IPs to check in one request. |
| max_servers | int | `100` | The maximum amount of servers allowed on one IP. |
| timeout | int | `30` | The request timeout. |

<details>
    <summary>Example(s)</summary>

If we want to check for duplicate servers every `60` seconds (1 minute) with the max servers per IP being set to `1000` and the server limit per request set to `200`, we'd use the following.

```json
{
    "enabled": true,
    "interval": 60,
    "limit": 200,
    "max_servers": 1000,
    "timeout": 30
}
```
</details>

### Remove Timed Out Object
This object contains settings for removing (setting offline) timed out servers based off of the platform's server timeout.

| Name | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| enabled | bool | `false` | Whether to enable remove timed out servers checking. |
| interval | int | `120` | How often to check for timed out servers in seconds. |
| platform_ids | []int | `[]` | A list of platform IDs to perform the server timed out check for. If left empty, all platforms are checked. |
| timed_out_time | int | `3600` | The server's timeout time. |
| timeout | int | `30` | The request timeout. |

<details>
    <summary>Example(s)</summary>

If we want to check for timed out servers for Counter-Strike: Source (ID `4`) every `120` seconds with a server timeout of `3600` seconds (an hour), we'd use the following.

```json
{
    "enabled": true,
    "interval": 120,
    "platform_ids": [4],
    "timed_out_time": 3600,
    "timeout": 30
}
```
</details>

## Credits
* [Christian Deacon](https://github.com/gamemann)