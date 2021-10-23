package dskstore

import (
    "crypto/sha1"

    "encoding/hex"
    "errors"

    "fmt"

    "io"
    "io/fs"

    "os"

    "path/filepath"

    "strings"
)

import (
    "github.com/rogpeppe/go-internal/lockedfile"
)

const (
    MaxPartitions = 16
    MaxLevels     = sha1.Size * 2
)

const (
    dirMode  = 0777
    fileMode = 0666
)

type DskStore struct {
    path string
    prt  int
    lvl  int
}

func NewDskStore(path string, prt, lvl int) (d *DskStore, err error) {
    if prt > MaxPartitions {
        err = fmt.Errorf("Too many partitions requested (%v)", prt)
        return
    }

    if lvl > MaxLevels {
        err = fmt.Errorf("Too many cache levels requested (%v)", lvl)
        return
    }

    if path, err = filepath.Abs(path); err != nil {
        return
    }

    d = &DskStore{
        path,
        prt,
        lvl,
    }

    if err = d.createPartitions(); err != nil {
        return
    }

    return
}

func (d *DskStore) Exists(fn string) (exists bool, err error) {
    base, name := d.getCachePath(fn)

    inf, err := os.Stat(filepath.Join(base, name))
    if err != nil {
        if errors.Is(err, fs.ErrNotExist) {
            err = nil
        }
        return
    }

    if !inf.Mode().IsRegular() {
        err = fmt.Errorf("object is not a file: %v", fn)
        return
    }

    exists = true

    return
}

func (d *DskStore) Store(fn string, src io.Reader) (err error) {
    base, name := d.getCachePath(fn)
    if err = os.MkdirAll(base, dirMode); err != nil {
        return
    }

    if err = lockedfile.Write(filepath.Join(base, name), src, fileMode); err != nil {
        return
    }

    return
}

func (d *DskStore) Retrieve(fn string) (data []byte, err error) {
    base, name := d.getCachePath(fn)

    if data, err = lockedfile.Read(filepath.Join(base, name)); err != nil {
        if errors.Is(err, fs.ErrNotExist) {
            err = fmt.Errorf("object not found: %v", fn)
        }
        return
    }

    return
}

func (d *DskStore) Clean() (err error) {
    if err = os.RemoveAll(d.path); err != nil {
        return
    }

    if err = d.createPartitions(); err != nil {
        return
    }

    return
}

func (d *DskStore) createPartitions() (err error) {
    z := fmt.Sprintf("%%0%dd", len(fmt.Sprint(d.prt)))

    for i := 0; i < d.prt; i++ {
        p := filepath.Join(d.path, fmt.Sprintf(z, i))
        if err = os.MkdirAll(p, dirMode); err != nil {
            return
        }
    }

    return
}

func (d *DskStore) getCachePath(fn string) (base string, name string) {
    bs := sha1.Sum([]byte(fn))

    base = filepath.Join(d.path, fmt.Sprintf("%x", int(bs[0] >> 4) % d.prt))
    for i := 0; i < d.lvl; i++ {
        base = filepath.Join(base, fmt.Sprintf("%x", (bs[i / 2] >> (-4 * (i % 2 - 1))) & 0xf))
    }

    name = hex.EncodeToString(bs[:]) + filepath.Ext(strings.TrimLeft(fn, "."))

    return
}
