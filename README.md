# deltadiff

This is a delta diffing library and CLI written in Go. It implements a basic rolling hash algorithm, as well as non-rolling hash, block based delta diffing.

I must emphasize how **experimental** this is.

> Note that this could fail on large files. Also the code is using the machine dependent `int` type and I only tested it on a 64-bit machine.

# What is this for

You have two similar files, `target` and `base`. You want to calculate the difference between them and apply the required changes in order to transform `base` into `target`. This is easy if both files are located in the same machine, but what if they are not? Then, this is how you do it:

First, you calculate the signature of `base`; then, you transfer the signature over the network to the machine where `target` is stored, and calculate the delta (difference) between signature and target. Finally, you transfer the delta over to the machine where `base` is located, and now you can apply the delta on to `base`, having it to become `target`.

```
signature = calculate-signature(base)
delta = calculate-delta(signature, target)
result = patch(base, delta)
```

Now `result` should now be the same as `target`.

# Signature

The signature is a byte encoded file in the following format:

```
+-----------+-------------+
|   size    |   content   |
+-----------+-------------+
| 2 bytes   | hasher code |
| 4 bytes   | block size  |
| 4 bytes   | base size   |
| remaining | blocks      |
+-----------+-------------+
```

**hasher code** This is the code that represents the hashing method, called a hasher, used when building the signature.

**block size** This is the size of the block that was used when building the signature.

**base size** This is the size of `base`. It's necessary because delta will only have access to the signature, so it doesn't know the size of `base` unless we tell it.

**blocks** The remaining of the signature is a sequence of `block size`-sized blocks.

# Delta

The delta is a byte encoded file containing only a sequence of operations. Those operations can be of one of two types: Read or Write. Read operations always refer to `base` and Write operations always refer to `target`. Therefore, Read operations can be read as "Read from base" and Write operations can be read as "Write from target".

For example, imagine the following `target` and `base`:

```
target = aaaabbbbccccddddeeee
base   = aaaabbbbccccddeeeeee
```

In order to transform `base` into `target`, this is a sequence of `read` and `write` operations that can be applied:

```
1. read:0-12
2. write:12-14
3. read:12-16
4. write:18-20
```

Where `read:0-12` means "read from `base`, from byte 0 to byte 12", and `write:12-14` means "write from `target`, from byte 12 to byte 14".

That said, for the read operations the delta file only contains the positional information of reads, and the actual content of writes. Therefore the delta file will actually look like this:

```
read:0-12
write:ccdd
read:12-16
write:ee
```

# Example (lib)

```go
package main

import (
	"bytes"
	"fmt"
	"github.com/xrash/deltadiff"
	"os"
)

var target = "aaaabbbbccccddddeeee"
var base = "aaaabbbbccccddeeeeee"

func main() {
	targetBuffer := bytes.NewBufferString(target)
	baseBuffer := bytes.NewBufferString(base)
	signatureBuffer := bytes.NewBuffer(nil)
	deltaBuffer := bytes.NewBuffer(nil)
	resultBuffer := bytes.NewBuffer(nil)

	sc := &deltadiff.SignatureConfig{
		Hasher:    "polyroll",
		BlockSize: 4,
		BaseSize:  len(base),
	}

	if err := deltadiff.Signature(baseBuffer, signatureBuffer, sc); err != nil {
		panic(err)
	}

	dc := &deltadiff.DeltaConfig{
		Debug:       true,
		DebugWriter: os.Stderr,
	}

	if err := deltadiff.Delta(signatureBuffer, targetBuffer, deltaBuffer, dc); err != nil {
		panic(err)
	}

	baseBuffer = bytes.NewBufferString(base)
	if err := deltadiff.Patch(baseBuffer, deltaBuffer, resultBuffer); err != nil {
		panic(err)
	}

	fmt.Println(target == resultBuffer.String())
	fmt.Println(target, resultBuffer.String())
}
```

This will output:


```
match	0:0-4
match	1:4-8
match	2:8-12
match	3:14-18
op	read:0-12
op	write:12-14
op	read:12-16
op	write:18-20
true
aaaabbbbccccddddeeee aaaabbbbccccddddeeee
```

# Example (CLI)

Run `deltadiff signature <base> <signature>` to calculate the signature:

```
$ deltadiff signature myfile-modified.jpg signature
```

Run `deltadiff delta <signature> <target> <delta>` to calculate the delta:

```
$ deltadiff delta delta signature myfile.jpg delta
```

Run `deltadiff patch <base> <delta> <result>` to apply the patch:

```
$ deltadiff patch myfile-modified.jpg delta result.jpg
```

Now `result.jpg` is the same as `myfile.jpg`

# Signature and Delta options

Both the library and the CLI have some options you can tweak. 

Signature has the following configuration:

```go
type SignatureConfig struct {
	Hasher    string
	BlockSize int
	BaseSize  int
}
```

Those options can be set when using the CLI through `--hasher` and `--block-size`.

`Hasher` can be `md5`, `crc32` or `polyroll`. The default value is `polyroll` - a custom, experimental rolling hash algorithm.

`BlockSize` defaults to 1024.

Delta has the following configuration:

```go
type DeltaConfig struct {
	Debug       bool
	DebugWriter io.Writer
}
```

Debugging can be turned on in the CLI through `--debug` and `--debug-file`. When set, it outputs the block matches and the sequence of operations.

# Installing the CLI

Run the command below:

```
$ go install github.com/xrash/deltadiff/cmd/deltadiff@latest
```

Now you should be able to run:

```
$ deltadiff
```

# Compiling CLI locally

Checkout and run:

```
$ make
```

Now move `./bin/deltadiff` wherever you want.

# Running tests


Checkout and run:

```
$ make test
```

It will take a while.
