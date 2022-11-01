package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/filecoin-project/venus/venus-devtool/util"
	"github.com/filecoin-project/venus/venus-shared/actors"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "state-type-gen",
		Usage:                "generate types related codes for go-state-types",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "dst",
			},
		},
		Action: run,
	}

	app.Setup()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %v\n", err) // nolint: errcheck
	}
}

var prePath = "github.com/filecoin-project/go-state-types/builtin"

type pendingPkg struct {
	name string
	path string
	ver  actors.Version
}

var pendingPkgs = func() map[string]*pendingPkg {
	pkgs := make(map[string]*pendingPkg, 4)
	list := []string{"market", "miner", "verifreg"}
	pkgs["paych"] = &pendingPkg{
		name: "paych",
		ver:  actors.Version8,
		path: fmt.Sprintf("%s/v%v/%s", prePath, actors.Version8, "paych"),
	}
	for _, pkgName := range list {
		pkgs[pkgName] = &pendingPkg{
			name: pkgName,
			ver:  actors.Version(actors.LatestVersion),
			path: fmt.Sprintf("%s/v%v/%s", prePath, actors.LatestVersion, pkgName),
		}
	}

	return pkgs
}()

var (
	rootDir = fmt.Sprintf("github.com/filecoin-project/go-state-types/builtin/v%v/", actors.LatestVersion)
	skips   = map[string]struct{}{
		"State":          {},
		"MinerInfo":      {},
		"ConstructState": {},
		"Partition":      {},
		"Deadline":       {},
	}
	skipFuncs = map[string]struct{}{
		"ConstructState": {},
	}
	alias = map[string][]struct {
		pkgName string
		newName string
	}{
		"WithdrawBalanceParams": {
			{pkgName: "market", newName: "MarketWithdrawBalanceParams"},
			{pkgName: "miner", newName: "MinerWithdrawBalanceParams"},
		},
	}
	expectVals = map[string]struct{}{
		"NoAllocationID": {},
	}
)

func run(cctx *cli.Context) error {
	metas := make([]*metaVisitor, 0, len(pendingPkgs))
	for _, pkg := range toList(pendingPkgs) {
		location, err := util.FindPackageLocation(pkg.path)
		if err != nil {
			return err
		}

		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, location, filter, parser.AllErrors|parser.ParseComments)
		if err != nil {
			return err
		}

		visitor := &metaVisitor{
			pkgName: pkg.name,
		}
		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				ast.Walk(visitor, file)
			}
		}

		sort.Slice(visitor.f, func(i, j int) bool {
			return visitor.f[i] < visitor.f[j]
		})
		sort.Slice(visitor.t, func(i, j int) bool {
			return visitor.t[i] < visitor.t[j]
		})
		metas = append(metas, visitor)
	}

	return writeFile(cctx.String("dst"), metas)
}

func toList(pkgs map[string]*pendingPkg) []*pendingPkg {
	list := make([]*pendingPkg, 0, len(pkgs))
	for _, pkg := range pkgs {
		list = append(list, pkg)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].name < list[j].name
	})

	return list
}

func filter(fi fs.FileInfo) bool {
	if strings.Contains(fi.Name(), "cbor_gen.go") {
		return false
	}
	if strings.Contains(fi.Name(), "_test.go") {
		return false
	}
	if strings.Contains(fi.Name(), "invariants.go") {
		return false
	}
	if strings.Contains(fi.Name(), "methods.go") {
		return false
	}
	return true
}

type metaVisitor struct {
	pkgName string
	f       []string // function
	t       []string // type
	v       []string // value | const
}

func (v *metaVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if st, ok := node.(*ast.TypeSpec); ok {
		if !st.Name.IsExported() {
			return v
		}

		name := st.Name.Name
		_, ok = st.Type.(*ast.StructType)
		_, ok3 := st.Type.(*ast.Ident)
		if _, ok2 := skips[name]; !ok2 && (ok || ok3) {
			v.t = append(v.t, name)
		}
	} else if ft, ok := node.(*ast.FuncDecl); ok {
		if !ft.Name.IsExported() || ft.Recv != nil {
			return v
		}

		name := ft.Name.Name
		if _, ok := skipFuncs[name]; !ok {
			v.f = append(v.f, name)
		}
	} else if vt, ok := node.(*ast.ValueSpec); ok {
		if !vt.Names[0].IsExported() || len(vt.Names) == 0 {
			return v
		}

		if _, ok := expectVals[vt.Names[0].Name]; ok {
			v.v = append(v.v, vt.Names[0].Name)
		}
	}

	return v
}

func writeFile(dst string, metas []*metaVisitor) error {
	var fileBuffer bytes.Buffer
	fmt.Fprintf(&fileBuffer, "// Code generated by github.com/filecoin-project/venus/venus-devtool/state-type-gen. DO NOT EDIT.\npackage %s\n\n", "types")

	// write import
	fmt.Fprintln(&fileBuffer, "import (")
	for _, meta := range metas {
		fmt.Fprintf(&fileBuffer, "\"%v\"\n", pendingPkgs[meta.pkgName].path)
	}
	fmt.Fprintln(&fileBuffer, ")\n")

	for _, meta := range metas {
		fmt.Fprintf(&fileBuffer, "////////// %s //////////\n", meta.pkgName)
		for _, typ := range meta.t {
			if vals, ok := alias[typ]; ok {
				for _, val := range vals {
					if val.pkgName == meta.pkgName {
						fmt.Fprintf(&fileBuffer, "type %s = %s.%s\n", val.newName, meta.pkgName, typ)
					}
				}
			} else {
				fmt.Fprintf(&fileBuffer, "type %s = %s.%s\n", typ, meta.pkgName, typ)
			}
		}

		for _, f := range meta.f {
			fmt.Fprintf(&fileBuffer, "var %s = %s.%s\n", f, meta.pkgName, f)
		}

		for _, v := range meta.v {
			fmt.Fprintf(&fileBuffer, "const %s = %s.%s\n", v, meta.pkgName, v)
		}
		fmt.Fprintln(&fileBuffer, "\n")
	}

	formatedBuf, err := util.FmtFile("", fileBuffer.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(dst, formatedBuf, 0o755)
}