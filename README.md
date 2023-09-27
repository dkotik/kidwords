# Kid Words, Alpha

Data encoding accessible to children for generating passwords and paper keys or splitting them into shards using Shamir's Secret Sharing algorithm.

## Release Checklist

- [ ] Add Shamir's Secret Sharing key splitting.
- [ ] Add Shamir's Secret Sharing key re-combination.
- [ ] Harden Shamir's Secret Sharing algorithm with `mod Prime`.
  - See https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing
  - Another alternative implementation uses prime
  - Prime should be configurable
- [ ] Add HTML SeparatorFunc.

## Command Line Tool

```sh
# Command line tool installation:
go install github.com/dkotik/kidwords/cmd/kidwords@latest
kidwords --help
```

### Key Splitting

The secret is compressed using Zstd algorithm before getting split into eight shards. Quorum is set using `--quorum=3` flag. The number of shards is limited to eight in order to use additional 13 bites for an error detection code. The shard ordinal and the error detection code are expressed as two additional words appended to the end of each shard.

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
