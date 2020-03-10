CLInium - BDD style tests for command line tools
================================================

Go code (golang) that provides tools for testing that your command line tool will behave as you intend.

Features include:
  * Check if command was successful or failed
  * Assertions on stdout
  * Execution of interactive commands

See it in action:
```go
package yours

import (
  "testing"
  . "github.com/SteffiPeTaffy/clinium"
)

func Test_YourAwesomeTest(t *testing.T) {
	myCli := NewCli(t, "../foo/bar/myCli")

    myCli.Expect("foo", "--bar", "baz").
      ToSucceed("This should work").
      ToNotContain("banana")
}
```

------


Supported go versions
==================

We support the three major Go versions, which are 1.11, 1.12, and 1.13 at the moment.

------

Contributing
============

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue. Extra credit for those using CLInium to write the test code that demonstrates it.

------

License
=======

This project is licensed under the terms of the MIT license.