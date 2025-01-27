package provider

import (
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/provider/custom"
	"github.com/metafates/mangal/provider/generic"
	"github.com/metafates/mangal/provider/mangadex"
	"github.com/metafates/mangal/provider/manganato"
	"github.com/metafates/mangal/provider/manganelo"
	"github.com/metafates/mangal/provider/mangapill"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

func init() {
	for _, conf := range []*generic.Configuration{
		manganelo.Config,
		manganato.Config,
		mangapill.Config,
	} {
		conf := conf
		builtinProviders = append(builtinProviders, &Provider{
			ID:   conf.ID(),
			Name: conf.Name,
			CreateSource: func() (source.Source, error) {
				return generic.New(conf), nil
			},
		})
	}
}

type Provider struct {
	ID           string
	Name         string
	CreateSource func() (source.Source, error)
}

func (p Provider) String() string {
	return p.Name
}

const CustomProviderExtension = ".lua"

var builtinProviders = []*Provider{
	{
		ID:   mangadex.ID,
		Name: mangadex.Name,
		CreateSource: func() (source.Source, error) {
			return mangadex.New(), nil
		},
	},
}

func Builtins() []*Provider {
	return builtinProviders
}

func Customs() []*Provider {
	files, err := filesystem.Api().ReadDir(where.Sources())

	if err != nil {
		return make([]*Provider, 0)
	}

	paths := lo.FilterMap(files, func(f os.FileInfo, _ int) (string, bool) {
		if filepath.Ext(f.Name()) == CustomProviderExtension {
			return filepath.Join(where.Sources(), f.Name()), true
		}
		return "", false
	})
	providers := make([]*Provider, len(paths))

	for i, path := range paths {
		name := util.FileStem(path)
		path := path
		providers[i] = &Provider{
			ID:   custom.IDfromName(name),
			Name: name,
			CreateSource: func() (source.Source, error) {
				return custom.LoadSource(path, true)
			},
		}
	}

	return providers
}

func Get(name string) (*Provider, bool) {
	for _, provider := range Builtins() {
		if provider.Name == name {
			return provider, true
		}
	}

	for _, provider := range Customs() {
		if provider.Name == name {
			return provider, true
		}
	}

	return nil, false
}
