---
title: "Kid Words, Beta"
menu:
  main:
    weight: 10
---

**Durable and accessible paper key encoding that children can use.**

> **Warning:** beta version is not stable and subject to iteration!

Printable paper keys are occasionally used as the last resort for recovering account access. They increase security by empowering a user to wrestle control of a compromised account from an attacker.

Most paper keys are encoded using BIP39 convention into a set of words. The final few words encode the integrity of the key with a cyclical redundancy check. When printed and stored, such keys are not durable because they can succumb to minor physical damage.

Kid Words encoding increases key durability by splitting the key using [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing) algorithm into shards and encoding each shard and its checksum into a group of four-letter English words.

## Explore the documentation

- [Benefits]({{< relref "benefits" >}})
- [Using as Library]({{< relref "library" >}})
- [Using as Command Line Tool]({{< relref "command-line" >}})
- [Specification]({{< relref "specification" >}})
- [Planned Features]({{< relref "planned-features" >}})
