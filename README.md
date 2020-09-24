<p align="center">
  <img src="https://github.com/talhof8/gasper/blob/master/assets/logo.png?raw=true" alt="Gasper Logo"/>
</p>

# Gasper

![](https://img.shields.io/github/issues/talhof8/gasper)
![](https://img.shields.io/github/stars/talhof8/gasper)
![](https://img.shields.io/github/license/talhof8/gasper)
![](https://img.shields.io/twitter/url?url=https%3A%2F%2Fgithub.com%2Ftalhof8%2Fgasper)

Back-up your files in a distributed manner, across multiple stores of your choice, by splitting them to shares. 
Retrieve them at any point, with only a minimum number of shares required for retrieval.

Gasper is based on the awesome [Shamir's Secret Sharing algorithm](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing). 

## Supported stores

| Type              | Description           | Attributes                |
| ----------------- |-----------------------| --------------------------|
| `local`      | Store share in a local directory | `directory-path` (string) |

Feel free to contribute your own stores - S3, Google Drive, Twitter, FTP, or anything else you'd like :)

## Installation
```
go get -u github.com/talhof8/gasper
```

# Demo
#### Using local store
![](assets/demo-local.gif)


## Usage
#### Store
```
gasper store --stores-config </path/to/stores.json> --file <file> [--encrypt --salt <valid-aes-salt> --share-count <count> --shares-threshold <min-threshold> --verbose]
```
Outputs file ID and checksum on success which should be used for retrieval.

#### Retrieve
```
gasper retrieve --stores-config </path/to/stores.json> --file-id <file-id> --destination <some-destination> [--checksum <some-checksum> --encrypt --salt <valid-aes-salt> --verbose]
```

#### Delete
Best effort deletion.
```
gasper delete --stores-config </path/to/stores.json> --file-id <file-id> [--verbose]
```

Stores configuration file:
```
{
  "stores": [
    {
      "type": "<type>",
      "<attribute>": "<value>",
      "<attribute>": "<value>",
      "<attribute>": "<value>"
    },
    {
      "type": "<type>",
      "<attribute>": "<value>",
      ...
    }
  ]
}
```

## License
Gasper is released under GPL. See LICENSE.txt.
