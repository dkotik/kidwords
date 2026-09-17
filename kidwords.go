/*
Package kidwords provides durable and accessible paper key encoding that children can use.

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

## Inspired By

- [Horcrux][horcrux]

horcrux: https://github.com/jesseduffield/horcrux/tree/master
*/
package kidwords

import (
	"bytes"
	"io"
)

type Codec interface {
	EncodeBytes([]byte) []byte
	DecodeBytes([]byte) ([]byte, error)
	EncodeStream(io.Writer, io.Reader) error
	DecodeStream(io.Writer, io.Reader, *bytes.Buffer) error
}

type kidwords struct {
	// Content Dictionary
	// Checksum Dictionary
}

func New() Codec {
	return nil
}
