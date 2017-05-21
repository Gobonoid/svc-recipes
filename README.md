# svc-recipes
Recipes contain a lot of information such as cuisine, customer ratings &amp; comments, stock levels and diet types.

### Requirements
* Go 1.8+
* Make
* Glide (https://github.com/Masterminds/glide)

### How to test, build and run

```
    make
    make run
```

or just using go commands

```
    go test -v --cover --race  `glide novendor`
    go run main.go

```

### Assumptions:
* No validation of incoming elements is required
* No database used - hence some weird work arounds in model
* Code doesn't have to be perfect and I dont't need to waste to much time

### General thoughts:
GW is incorrect name in terms of architecture and this rather should be service 