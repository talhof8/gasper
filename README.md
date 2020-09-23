# Gasper
Split and distribute your most sacred, sensitive files across multiple stores 
of your choice, with only a minimum number of them needed to retrieve them back.

Gasper is based on the awesome [Shamir's Secret Sharing algorithm](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing). 

## Supported stores

| Type              | Description           | Attributes                |
| ----------------- |-----------------------| --------------------------|
| `local`      | Store share in a local directory | `directory-path` (string) |

Feel free to open a Pull Request and add your own :)

## Installation
```
go get -u github.com/talhof8/gasper
```

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
Best effort shares deletion.
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