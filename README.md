# Qor Slug

Slug package is an extension for qor. It provides an easy way to create a pretty URL for your model.

## Usage

```go
import (
	"github.com/jinzhu/gorm"
	"github.com/qor/slug"
)

type Product struct {
	gorm.Model
	Name            string
	NameWithSlug    slug.Slug
}
```

## License

Released under the MIT License.
