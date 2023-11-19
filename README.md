# wgs84 v2 WIP

```sh
go get github.com/wroge/wgs84/v2
```

I am currently in the process of rewriting the package. Some things will change and some new features will be added. One of these features is the support of NTv2 grid transformations and other projections, such as Krovak. If you would like to help or have any comments, please report them in the issues.

### Web Mercator

```go
package main

import (
	"fmt"

	"github.com/wroge/wgs84/v2"
)

func main() {
	transform := wgs84.Transform(wgs84.EPSG(4326), wgs84.EPSG(3857)).Round(3)

	east, north, _ := transform(10, 50, 0)

	fmt.Println(east, north)
	// 1.113194908e+06 6.446275841e+06

	// echo 10 50 | cs2cs +init=epsg:4326 +to +init=epsg:3857 -d 3
	// 1113194.908     6446275.841
}
```

### OSGB

```go
package main

import (
	"fmt"

	"github.com/wroge/wgs84/v2"
)

func main() {
	transform := wgs84.Transform(wgs84.EPSG(4326), wgs84.EPSG(27700)).Round(3)

	east, north, h := transform(-2.25, 52.25, 0)

	fmt.Println(east, north, h)
	// 383029.296 261341.615 0

	// echo -2.25 52.25 | cs2cs +init=epsg:4326 +to +init=epsg:27700 -d 3
	// 383029.296 261341.615 0.000
}
```
