---
title: "Using as Library"
weight: 30
---

```go
// To install the library run shell command:
//
// $ go get github.com/dkotik/kidwords@latest

func main() {
  // break a secret key into shards
  shards, err := kidwords.Split(
    []byte("secret paper key"), // encoding target
    12,                          // number of shards
    4,                           // quorum number of shards
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
