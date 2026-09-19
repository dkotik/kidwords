---
title: "Using as Command Line Tool"
weight: 40
---

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
