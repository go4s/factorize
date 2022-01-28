package main

import (
    "embed"
    "flag"
    "html/template"
    "os"
    "path/filepath"
    "strings"
)

const OutputFiles = "boilerplate/*.go.tpl"

//go:embed boilerplate/*.go.tpl
var BoilerPlates embed.FS

var (
    pkgName string
    typName string
    dir     string

    pwd  string
    tpls *template.Template
    data map[string]interface{}
)

func init() {
    flag.StringVar(&dir, "o", ".", "-o=OutputDir")
    flag.StringVar(&pkgName, "p", "", "-p=PkgName")
    flag.StringVar(&typName, "t", "", "-t=TypName")
    flag.Parse()
    var err error
    if tpls, err = template.ParseFS(BoilerPlates, OutputFiles); err != nil {
        panic(err)
    }
    if err = os.MkdirAll(filepath.Join(pwd, dir, pkgName), 644); err != nil {
        panic(err)
    }
    if pwd, err = os.Getwd(); err != nil {
        panic(err)
    }
    data = map[string]interface{}{
        "PktName": pkgName, "TypName": strings.Title(typName),
        "BinVer": BinVer, "BinName": BinName,
        "BinYear": BinYear, "BinAuth": BinAuth,
    }
}

func main() {
    for _, t := range tpls.Templates() {
        name := strings.TrimSuffix(t.Name(), ".tpl")
        if strings.HasSuffix(name, "kind.go") {
            continue
        }
        outputFile, err := os.OpenFile(filepath.Join(pwd, dir, pkgName, name), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
        if err != nil {
            panic(err)
        }
        if err := t.Execute(outputFile, data); err != nil {
            panic(err)
        }
        _ = outputFile.Close()
    }
    outputFile, err := os.OpenFile(filepath.Join(pwd, dir, pkgName, pkgName+".go"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
    defer outputFile.Close()
    if err != nil {
        panic(err)
    }
    if err = tpls.ExecuteTemplate(outputFile, "kind.go.tpl", data); err != nil {
        panic(err)
    }
}
