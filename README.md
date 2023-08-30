### mitum-nft v2

*mitum-nft v2* is a nft contract model based on the second version of mitum(aka [mitum2](https://github.com/ProtoconNet/mitum2)).

#### Features,

* user defined contract account state policy: collection.
* collection: user defined nft collection.
* *LevelDB*, *Redis*: as mitum2 does, *LevelDB* and *Redis* can be primary storage.
* reference nft standard: ERC-721
* multiple collection policy for one contract account.

#### Installation

Before you build `mitum-nft`, make sure to run `docker run`.

```sh
$ git clone https://github.com/wyuinche/mitum-nft

$ cd mitum-nft

$ git checkout -t origin/refactor

$ go build -o ./mitum-nft
```

#### Run

```sh
$ ./mitum-nft init --design=<config file> <genesis file>

$ ./mitum-nft run <config file>
```

[standalong.yml](standalone.yml) is a sample of `config file`.

[genesis-design.yml](genesis-design.yml) is a sample of `genesis design file`.

To run a node with the default settings, enter the following command:

```sh
$ ./mitum-nft init --design=standalone.yml genesis-design.yml

$ ./mitum-nft run --design=standalone.yml --dev.allow-consensus
```

Docker and Mongodb must be installed to run with digest api.

```sh
$ docker run --name mnft -it -p 27017:27017 -d mongo
```

To run a node without digest api, remove the following from the [standalong.yml](standalone.yml) file.

```yml
digest:
  network:
    bind: http://localhost:54320
    url: http://localhost:54320
  database: 
    uri: mongodb://127.0.0.1:27017/mnft
```
