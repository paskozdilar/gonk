# gonk

Pure Go dynamic loader experiment.

This repository is written for learning purposes, please don't use it IRL.

## Current state

We're able to load an `int add(int, int)` C function from shared library into
memory, execute it, and get the return value back.

```
$ ./example/build.sh
$ go run .
Result from C: 42
```

## Roadmap

- [x] Dynamic Symbol Loading
- [ ] Data Relocation
- [ ] Global Offset Table & Procedure Linkage Table
- [ ] Fine-Grained Memory Permissions
- [ ] BSS Section Initialization
