---
title: "Specification"
weight: 50
---

1. The secret is split into data shards encoded with four-letter English words:

   ```go
   type Shard struct {
    Index    uint8  // decoded from (line_number_prefix - 1)
    Data     []byte // decoded from English nouns
    Checksum []byte // decoded big-endian from four English verbs (32bit=4bytes)
   }
   ```

2. Each line of an encoded shard begins with a shard index number to simplify recognition of shard boundaries when reading or using optical character recognition algorithms. Shard index begins with **1**, never with zero.

3. Nouns and verbs of one shard can be mixed with each other, but they must always follow the exact left-to-right order relative to each other. Take two or less bytes of *Data* and one byte or none of *Checksum*. Repeat until both *Data* and *Checksum* bytes are fully encoded:

   ```text
   1 noun noun noun verb
   1 noun noun noun verb
   1 noun noun noun verb
   1 verb
   ```

4. Multiple shards can be encoded next to each other:

   ```text
   1 noun noun noun verb 2 noun noun noun verb 3 noun noun noun verb
   1 noun noun noun verb 2 noun noun noun verb 3 noun noun noun verb
   1 noun noun noun verb 2 noun noun noun verb 3 noun noun noun verb
   1 noun verb           2 noun verb           3 noun verb
   ```

5. Paper keys should be deterministically fingerprinted in order to avoid the expensive password hashing for each paper key, when a user authenticates using one of the keys. Take the last byte of every `Shard.Data`. It contains a random `x` value with which the resulting polynomial `y` value is computed.

   A group of shards will have a rare combination of `.Index` to `x` that can be matched with a secret without performing any cryptographic operations on the secret. This combination does not reveal anything about the secret.

   Fingerprint values should have a unique constraint for each user.
