# Kid Words, _Alpha_

Provides durable and accessible paper key encoding that children can use.

Printable paper keys are occasionally used as the last resort for recovering account access. They increase security by empowering a user with the ability to wrestle control of a compromised account from an attacker.

Most paper keys are encoded using BIP39 convention into a set of words. The final few words encode the integrity of the key with a cyclical redundancy check. When printed and stored, such keys are not durable because they can be lost to minor physical damage.

Kid Words package or command line tool increases key durability by splitting the key using [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing) algorithm into shards and encoding each shard using a dictionary of 256 four-letter English nouns.

## Benefits

- Keys can be recovered from partially damaged paper.
- Shards can be transmitted and memorized by children.
- Shards are easier to speak over poor radio or telephone connection, which can save time during an emergency.
- Key shards can be hidden in several physical locations by cutting the paper into pieces. Once a configurable quorum of shards, three by default, is gathered back, the key can be restored.
- Shards can easily be obfuscated by sequencing:
  - toys or books on a shelf
  - pencil scribbles on paper
  - objects or signs in a Minecraft world
  - emojis
- Command line tool can apply all of the above benefits to:
  - important passwords to rarely accessed accounts that do not support paper keys
  - conventional BIP39 keys

## Release Checklist

- [ ] Add Shamir's Secret Sharing key re-combination.
- [ ] Harden Shamir's Secret Sharing algorithm with `mod Prime`.
  - See https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing
  - Another alternative implementation uses prime
  - Prime should be configurable
- [ ] Add HTML SeparatorFunc.
- [ ] Add Emoji dictionary

## Command Line Tool

```sh
# Command line tool installation:
go install github.com/dkotik/kidwords/cmd/kidwords@latest
kidwords --help
```

### Key Splitting

The secret is compressed using Zstd algorithm before getting split into eight shards. Quorum is set using `--quorum=3` flag.

> The number of shards is limited to eight in order to use additional 13 bites for an error detection code. The shard ordinal and the error detection code are expressed as two additional words appended to the end of each shard. (This is wrong - the last byte encodes a random shard ordinal from 0-255)

When the quorum is set to `3` any three of the shards will be sufficient to recover the secret. If the quorum is set to `8`, every single shard will be required.

## Library

```go
// In shell: $ go get github.com/dkotik/kidwords@latest

func main() {
  w, err := kidwords.NewWriter(os.Stdout)
  if err != nil {
    panic(err)
  }
  _, _ = w.Write([]byte("test")) // will output words  
}
```
