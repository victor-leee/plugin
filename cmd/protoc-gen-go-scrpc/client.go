package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
)

func generateClient(gen *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + "_go-scrpc.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-scrpc. DO NOT EDIT.")
	g.P("// Path ", file.Desc.Path())
	g.P("// This is the client side implementation")
	g.P("// If you would like to generate the server side, please add --go-scrpc_opt=server=true to your protoc command")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	g.P("var client ", scrpcPkg.Ident("Client"))
	g.P(`func init() {
client = scrpc.NewClient()
}
`)

	// generate services
	for _, service := range file.Services {
		generateService(g, service)
	}
	g.P()
	generateConfigCenter(g)
}

func generateService(g *protogen.GeneratedFile, serviceDef *protogen.Service) {
	// generate interface and default implementation based on serviceDef
	// generate interface first
	g.P("type ", serviceDef.GoName, " interface{")
	for _, method := range serviceDef.Methods {
		generateMethod(g, method)
	}
	g.P("}")

	// then we generate a default impl for the interface
	g.P()
	g.P("type ", serviceDef.GoName, "Impl", " struct{")
	g.P("}")
	for _, method := range serviceDef.Methods {
		generateMethodImpl(serviceDef.GoName, g, method)
	}
}

func generateMethod(g *protogen.GeneratedFile, methodDef *protogen.Method) {
	g.P(methodDef.GoName, " (ctx ", ctxPkg.Ident("Context"), ", req *", methodDef.Input.GoIdent.GoName, ") (*", methodDef.Output.GoIdent.GoName, ", error)")
}

func generateMethodImpl(service string, g *protogen.GeneratedFile, methodDef *protogen.Method) {
	g.P("func (*", service, "Impl) ", methodDef.GoName, "(ctx ", ctxPkg.Ident("Context"), ", req *", methodDef.Input.GoIdent.GoName, ") (*", methodDef.Output.GoIdent.GoName, ", error){")
	generateBody(g, methodDef)
	g.P("}")
}

func generateBody(g *protogen.GeneratedFile, methodDef *protogen.Method) {
	g.P(fmt.Sprintf(`resp := &%s{}
err := client.UnaryRPCRequest(&scrpc.RequestContext{
Ctx: ctx,
Req: req,
ReqService: "%s",
ReqMethod: "%s",
SenderService: "%s",
Resp: resp,
})`, methodDef.Output.GoIdent.GoName, convertGoPath2DNS(file.GoImportPath.String()), methodDef.GoName, convertGoPath2DNS(defaultCfg.Service)))
	g.P("return resp, err")
}
