package plugin

import (
	"context"
	"time"

	"github.com/bufbuild/protoplugin"
	"github.com/mfridman/protoc-gen-go-json/gen"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	defaultFilenameSuffix = ".pb.json.go"
)

func Handle(ctx context.Context, w *protoplugin.ResponseWriter, r *protoplugin.Request) error {
	p, err := protogen.Options{}.New(r.CodeGeneratorRequest())
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	opt, err := parseOptions(r.Parameter())
	if err != nil {
		return err
	}
	if err := generate(p, opt); err != nil {
		p.Error(err)
	}

	response := p.Response()
	w.AddCodeGeneratorResponseFiles(response.GetFile()...)
	w.SetError(response.GetError())
	w.SetFeatureProto3Optional()
	return nil
}

func generate(p *protogen.Plugin, opt *gen.Options) error {
	for _, f := range p.Files {
		if len(f.Messages) == 0 {
			continue
		}
		gf := p.NewGeneratedFile(f.GeneratedFilenamePrefix+defaultFilenameSuffix, f.GoImportPath)
		if err := gen.ApplyTemplate(gf, f, opt); err != nil {
			gf.Skip()
			return nil
		}
	}
	return nil
}
