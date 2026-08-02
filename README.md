# Kid Words, _Alpha_

Provides durable and accessible paper key encoding that children can use.

**Warning: alpha version is not stable and subject to iteration!**

Printable paper keys are occasionally used as the last resort for recovering account access. They increase security by empowering a user with the ability to wrestle control of a compromised account from an attacker.

Most paper keys are encoded using BIP39 convention into a set of words. The final few words encode the integrity of the key with a cyclical redundancy check. When printed and stored, such keys are not durable because they can be lost to minor physical damage.

Kid Words package or command line tool increases key durability by splitting the key using [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing) algorithm into shards and encoding each shard using a dictionary of 256 four-letter English nouns.

## Benefits

- Keys can be recovered from partially damaged paper.
- Shards can be transmitted and memorized by children.
- Shards are easier to speak over poor radio or telephone connection, which can save time during an emergency.
- Shards can be hidden in several physical locations by cutting the paper into pieces. Once a configurable quorum of shards, four by default, is gathered back, the key can be restored.
- Shards can easily be obfuscated by sequencing:
  - toys or books on a shelf
  - pencil scribbles on paper
  - objects or signs in a Minecraft world
  - emojis
- Command line tool can apply all of the above benefits to:
  - important passwords to rarely accessed accounts that do not support paper keys
  - conventional BIP39 keys

## Development Checklist

- [ ] Use a separate 256-word alphabet for the checksum bytes.
- [ ] Harden Shamir's Secret Sharing algorithm with `mod Prime`.
  - See https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing
  - Another alternative implementation uses prime
  - Prime should be configurable?
- [ ] Implement modular HTTP service using https://templ.guide, HTMX, and https://github.com/mazznoer/csscolorparser for OKLCH colors, Zombie SQLite C-Go-less driver
- [ ] finish Argon hashing
- [ ] finish SQL store
- [ ] add BIP39 converter
- [ ] add Mongo store
- [ ] Add Emoji dictionary
- [ ] Add random password generator

## Using as Library

```go

import (
  "fmt"
  "os"

  // To install the library run shell command:
  //
  // $ go get github.com/dkotik/kidwords@latest
  "github.com/dkotik/kidwords"
  "github.com/dkotik/kidwords/shamir"
)

func main() {
  // break a secret key into shards
  shards, err := kidwords.Split(
    []byte("secret paper key"), // encoding target
    12,                         // number of shards
    4,                          // quorum number of shards
                                // needed to recover the original
  )
  if err != nil {
    panic(err)
  }
  if _, err = shards.Grid(
    3,  // number of table columns
    18, // number of characters to wrap the text at
  ).Write(os.Stdout); err != nil {
    panic(err)
  }

  // reconstitute the key back using a quorum of four shards
  key, err := shamir.Combine(shards[0:4])
  if err != nil {
    panic(err)
  }
  fmt.Println(string(key))
  // Output: secret paper key
}
```

## Using as Command Line Tool

```sh
$ go install github.com/dkotik/kidwords/cmd/kidwords@latest
$ kidwords split somePaperKey
🔑 Pick any 4 shards:
┌──────────────╥──────────────╥──────────────┐
│farm line belt║line hall cash║view home shot│
│beer crab pity║trap loot site║room turn tale│
│hour fund fuel║head flag pool║bank wind deal│
╞══════════════╬══════════════╬══════════════╡
│line hall cash║view home shot║help dirt turn│
│trap loot site║room turn tale║goat coat heir│
│head flag pool║bank wind deal║moss iron tour│
╞══════════════╬══════════════╬══════════════╡
│view home shot║help dirt turn║golf tape font│
│room turn tale║goat coat heir║pear debt dust│
│bank wind deal║moss iron tour║lake urge bush│
╞══════════════╬══════════════╬══════════════╡
│help dirt turn║golf tape font║wish risk cold│
│goat coat heir║pear debt dust║trap room card│
│moss iron tour║lake urge bush║firm moon root│
└──────────────╨──────────────╨──────────────┘
$ go run github.com/dkotik/kidwords/cmd/kidwords@latest combine
```
