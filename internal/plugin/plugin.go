package plugin

import (
	"context"

	"github.com/bufbuild/protoplugin"
	"github.com/mfridman/protoc-gen-go-json/internal/gen"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	defaultFilenameSuffix = ".pb.json.go"
)

func Handle(
	ctx context.Context,
	env protoplugin.PluginEnv,
	w protoplugin.ResponseWriter,
	r protoplugin.Request,
) error {
	p, err := protogen.Options{}.New(r.CodeGeneratorRequest())
	if err != nil {
		return err
	}
	opt, err := parseOptions(r.Parameter())
	if err != nil {
		return err
	}
	if err := generate(p, opt); err != nil {
		p.Error(err)
	}

	response := p.Response()
	w.AddCodeGeneratorResponseFiles(response.GetFile()...)
	w.AddError(response.GetError())
	w.SetFeatureProto3Optional()
	return nil
}

func generate(p *protogen.Plugin, opt *gen.Options) error {
	for _, f := range p.Files {
		if !f.Generate || len(f.Messages) == 0 {
			continue
		}
		g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+defaultFilenameSuffix, f.GoImportPath)
		if err := gen.ApplyTemplate(g, f, opt); err != nil {
			g.Skip()
			return nil
		}
	}
	return nil
}
