# dskstore
Store file in a multi-level disk structure

## Code example

```
d, err := dskstore.NewDskStore("storage", 6, 4)
if err != nil {
    panic(err)
}
```
- create storage with 6 partitions and 4-level depth

```
if err = d.Clean(); err != nil {
    panic(err)
}
```
- cleanup storage

```
if err = d.Store("00000001.txt", strings.NewReader("test data")); err != nil {
    panic(err)
}
```
- store provided content of "00000001.txt" in storage/3/9/9/a/f/99af904a8245824376beaab015e49a9d8a2db278.txt

```
if !d.Exists("00000001.txt") {
    panic("object doesn't exist")
}
```
- check object existence

```
data, err := d.Retrieve("00000001.txt")
if err != nil {
    panic(err)
}
```
- retrieve content of "00000001.txt"

## License

This project is licensed under the MIT License - see the
[LICENSE](https://github.com/pietro-a/dskstore/blob/master/LICENSE)
file for details.

