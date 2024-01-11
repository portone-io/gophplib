gophplib
========
[![doc icon]][doc link] [![coverage icon]][coverage link] [![report icon]][report link]

gophplib is a collection of PHP functions implemented in Go. We aim to achieve
100% byte-to-byte bug-to-bug behavioral equivalence with the original PHP
functions.

Currently we do not intend to implement all PHP functions, but only those that
are useful for our projects.

###### References
The following works were consulted but not used because they do not provide
behavior that is 100% compatible with the PHP version.

- https://github.com/hyperjiang/php
- https://github.com/syyongx/php2go

&nbsp;

--------
*gophplib* is primarily distributed under the terms of both the [MIT license]
and the [Apache License (Version 2.0)]. See [COPYRIGHT] for details.

[doc icon]: https://pkg.go.dev/badge/github.com/portone-io/gophplib.svg
[doc link]: https://pkg.go.dev/github.com/portone-io/gophplib
[coverage icon]: https://github.com/portone-io/gophplib/wiki/coverage.svg
[coverage link]: https://raw.githack.com/wiki/portone-io/gophplib/coverage.html
[report icon]: https://goreportcard.com/badge/github.com/portone-io/gophplib
[report link]: https://goreportcard.com/report/github.com/portone-io/gophplib
[MIT license]: LICENSE-MIT
[Apache License (Version 2.0)]: LICENSE-APACHE
[COPYRIGHT]: COPYRIGHT
