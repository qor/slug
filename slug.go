package slug

import (
	"database/sql/driver"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/qor/qor"
	"github.com/qor/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/qor/validations"
)

type Slug struct {
	Slug string
}

func (slug *Slug) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		slug.Slug = string(bytes)
	} else if str, ok := value.(string); ok {
		slug.Slug = str
	} else if strs, ok := value.([]string); ok {
		slug.Slug = strs[0]
	}
	return nil
}

func (slug Slug) Value() (driver.Value, error) {
	return slug.Slug, nil
}

var injected bool

func (Slug) ConfigureQorResource(res *admin.Resource) {
	Admin := res.GetAdmin()
	scope := Admin.Config.DB.NewScope(res.Value)

	if !injected {
		injected = true
		for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
			admin.RegisterViewPath(path.Join(gopath, "src/github.com/qor/slug/views"))
		}
		res.UseTheme("slug")
	}

	for _, field := range scope.Fields() {
		if field.Struct.Type == reflect.TypeOf(Slug{}) {
			name := strings.TrimSuffix(field.Name, "WithSlug")

			if meta := res.GetMeta(name); meta != nil {
				meta.Type = "slug"
			} else {
				res.Meta(&admin.Meta{Name: name, Type: "slug"})
			}
			if slugMeta := res.GetMeta(field.Name); slugMeta == nil {
				res.Meta(&admin.Meta{Name: field.Name, Type: "string"})
			}

			var fieldName = field.Name
			res.AddValidator(func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
				slug := utils.ToString(metaValues.Get(fieldName).Value)
				if slug == "" {
					return validations.NewError(record, fieldName, name+"'s slug can't be blank")
				} else if strings.Contains(slug, " ") {
					return validations.NewError(record, fieldName, name+"'s slug can't contains blank string")
				}
				return nil
			})

			res.IndexAttrs(append(res.IndexAttrs(), "-"+field.Name)...)
			res.ShowAttrs(append(res.ShowAttrs(), "-"+field.Name)...)
			res.EditAttrs(append(res.EditAttrs(), "-"+field.Name)...)
			res.NewAttrs(append(res.NewAttrs(), "-"+field.Name)...)
		}
	}
}
