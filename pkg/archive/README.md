# archive

This package provides utilities for working with archive formats.

## Functions

### Unzip

```
func Unzip(src, dest string) error
```

Extracts a ZIP archive from `src` to the `dest` directory. All files and folders in the archive will be extracted, preserving the directory structure.

#### Example

```go
import "github.com/Merith-TK/utils/pkg/archive"

err := archive.Unzip("example.zip", "outputDir")
if err != nil {
    // handle error
}
``` 