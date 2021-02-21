package rootloghandler

import (
	"embed"
	"io/fs"

	"github.com/blbgo/httpserver"
)

//go:embed templates
var templatesFS embed.FS

type templates struct{}

// NewTemplateFS provides the templates
func NewTemplateFS() httpserver.TemplateFSProvider {
	return templates{}
}

func (r templates) TemplateFS() (fs.FS, error) {
	subFS, err := fs.Sub(templatesFS, "templates")
	if err != nil {
		return nil, err
	}
	return subFS, nil
}
