# Kid Words, _Beta_

Durable and accessible paper key encoding that children can use.

**Warning: beta version is not stable and subject to iteration!**

Printable paper keys are occasionally used as the last resort for recovering account access. They increase security by empowering a user with the ability to wrestle control of a compromised account from an attacker.

Most paper keys are encoded using BIP39 convention into a set of words. The final few words encode the integrity of the key with a cyclical redundancy check. When printed and stored, such keys are not durable because they can succumb to minor physical damage.

Kid Words encoding increases key durability by splitting the key using [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing) algorithm into shards and encoding each shard and its checksum into a group of four-letter English words.

<details>
  <summary>Planned features for <strong>v1.0.0</strong> release. ↩</summary>

- [ ] Use a separate 256-verb alphabet for the checksum bytes.
- [ ] Add random password generator
- [ ] Add PDF generation
- [ ] Harden Shamir's Secret Sharing algorithm with `mod Prime`.
  - See https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing
  - Another alternative implementation uses prime
  - Prime should be configurable?
- [ ] Implement modular HTTP service using https://templ.guide, HTMX, and https://github.com/mazznoer/csscolorparser for OKLCH colors, Zombie SQLite C-Go-less driver
- [ ] finish Argon hashing
- [ ] finish SQL store
- [ ] Generate examples
</details>

## Benefits

The key shards can be recovered even by a child in emergency circumstances, when the paper had been partially damaged or the key shards are dictated over unstable communication medium.

- Shards are easier to speak over poor radio or telephone connection, which can save time during an emergency.
- Shards can be hidden in several physical locations by cutting the paper into pieces. Once a quorum of shards, four by default, is gathered, the key can be restored.
- Command line tool can apply all of the above benefits to:
  - important passwords to rarely accessed accounts that do not support paper keys
  - conventional BIP39 keys

## Using as Library

```go
// To install the library run shell command:
//
// $ go get github.com/dkotik/kidwords@latest
  
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

# Specification

1. The secret is split into data shards encoded with four-letter
English words:

  ```go
  type Shard struct {
  	Index    uint8  // decoded from (line_number_prefix - 1)
  	Data     []byte // decoded from English nouns
  	Checksum int32  // decoded big-endian from four English verbs
  }
  ```
2. Each line of an encoded shard begins with a shard index number to
simplify recognition of shard boundaries when reading or using optical
character recognition algorithms. Shard index begins with "**1**", never
with zero.
3. Nouns and verbs of one shard can be mixed with each other, but
they must always follow the exact left-to-right order relative to each
other. Date two or less bytes of _Data_ and one byte or none of 
_Checksum_. Repeat until both _Data_ and _Checksum_ bytes are fully
encoded:

  ```
  1 noun noun noun verb
  1 noun noun noun verb
  1 noun noun noun verb
  1 verb
  ```

4. Multiple shards can be encoded next to each other:

  ```
  1 noun noun noun verb 2 noun noun noun verb 3 noun noun noun verb
  1 noun noun noun verb 2 noun noun noun verb 3 noun noun noun verb
  1 noun noun noun verb 2 noun noun noun verb 3 noun noun noun verb
  1 noun verb           2 noun verb           3 noun verb
  ```
