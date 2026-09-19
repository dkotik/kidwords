---
title: "Planned Features for v1.0.0"
weight: 60
---

- [ ] Add random password generator
- [ ] Add PDF generation
- [ ] Harden Shamir's Secret Sharing algorithm with `mod Prime`.
  - See [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing)
  - Another alternative implementation uses prime
  - Prime should be configurable?
- [ ] Add Postgres SQL store
- [ ] The last byte of each share is a random `uint8` that identifies that share. That number could be used instead of the share index and would have 1 byte per shard.
