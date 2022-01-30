# Ukuleleweb

Ukuleleweb is a simple Wiki implementation in the style of the original C2 wiki.

## Few dependencies

The only dependencies are a markdown renderer and [Peter Bourgon's
`diskv` library](https://github.com/peterbourgon/diskv).

## Simple Storage

Wiki pages are simply stored in a directory, where file names are page
names and each file contains the page's markdown source. This way, you
can never lose your data, and it's also easy to put it under version
control.

## Run it

To run Ukuleleweb, create an empty directory where it can store data,
and run it using:

```
go run cmd/ukuleleweb/main.go --store_dir=/some/empty/directory
```
