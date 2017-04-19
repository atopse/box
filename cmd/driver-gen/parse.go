package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var buffer = bytes.NewBuffer(nil)

// Parse 解析Driver信息，并生成对应的解释文件
func Parse(pkgRealpath string) error {
	fileSet := token.NewFileSet()

	astPkgs, err := parser.ParseDir(fileSet, pkgRealpath, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
	}, parser.ParseComments)
	if err != nil {
		return err
	}

	drivers := make(map[string]map[string]interface{})

	for _, pkg := range astPkgs {
		for _, fl := range pkg.Files {
			for _, d := range fl.Decls {
				specDecl, ok := d.(*ast.GenDecl)
				if !ok || specDecl.Tok != token.TYPE {
					continue
				}
				structDecl, ok := specDecl.Specs[0].(*ast.TypeSpec)
				if !ok {
					continue
				}
				items, err := parserStructComments(specDecl.Doc)
				if err != nil {
					return err
				}
				//忽略不出来的项
				if len(items) == 0 {
					continue
				}
				structName := structDecl.Name.String()
				drivers[structName] = make(map[string]interface{})
				drivers[structName]["actions"] = []map[string]interface{}{}
				drivers[structName]["structName"] = structName
				drivers[structName]["pkg"] = pkg.Name
				for k, v := range items {
					drivers[structName][k] = v
				}
			}
			for _, d := range fl.Decls {
				specDecl, ok := d.(*ast.FuncDecl)
				if !ok || specDecl.Recv == nil {
					continue
				}
				exp, ok := specDecl.Recv.List[0].Type.(*ast.StarExpr) // Check that the type is correct first beforing throwing to parser
				if !ok || specDecl.Doc == nil || len(specDecl.Doc.List) == 0 {
					continue
				}
				structName := fmt.Sprint(exp.X)
				if _, ok := drivers[structName]; !ok {
					continue
				}
				data, err := parserComments(specDecl.Doc)
				if err != nil {
					return err
				}
				if len(data) == 0 {
					continue
				}
				drivers[structName]["actions"] = append(drivers[structName]["actions"].([]map[string]interface{}), data)
			}
		}
	}
	var buf bytes.Buffer
	for driver, data := range drivers {
		buf.Reset()
		data["now"] = time.Now()
		err = template.Must(template.New(driver).Parse(fileTpl)).Execute(&buf, data)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(fmt.Sprintf("gen_%s.go", driver), buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func parserStructComments(comments *ast.CommentGroup) (map[string]interface{}, error) {

	if comments == nil || len(comments.List) == 0 {
		return nil, nil
	}
	commentMap := make(map[string]interface{})
	for _, c := range comments.List {
		t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		if strings.HasPrefix(t, "@ignore") {
			return nil, nil
		}
		if !strings.HasPrefix(t, "@") {
			continue
		}
		kv := strings.SplitN(t, " ", 2)
		if len(kv) != 2 {
			return nil, errors.New("注释方式错误，参数和内容间需要用空格分割:" + t)
		}
		item := strings.Trim(strings.ToLower(kv[0]), "@")
		if item == "tags" {
			tags := []string{}
			for _, t := range strings.Split(kv[1], "#") {
				tags = append(tags, strings.TrimSpace(t))
			}
			commentMap["tags"] = tags
		} else if _, ok := commentMap[item]; !ok {
			commentMap[item] = strings.TrimSpace(kv[1])
		} else {
			return nil, errors.New("存在重复项 " + item)
		}
	}
	//如果无namespace则当做不处理
	if _, ok := commentMap["namespace"]; !ok {
		return nil, nil
	}
	if _, ok := commentMap["title"]; !ok {
		return nil, errors.New("@title不能为空")
	}
	return commentMap, nil
}

func parserComments(comments *ast.CommentGroup) (map[string]interface{}, error) {
	if comments == nil || len(comments.List) == 0 {
		return nil, nil
	}
	commentMap := make(map[string]string)
	for _, c := range comments.List {
		t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		if strings.HasPrefix(t, "@ignore") {
			return nil, nil
		}
		if !strings.HasPrefix(t, "@") {
			if _, ok := commentMap["@desc"]; !ok {
				commentMap["@desc"] = t
			}
			continue
		}
		kv := strings.SplitN(t, " ", 2)
		if len(kv) != 2 {
			return nil, errors.New("注释方式错误，参数和内容间需要用空格分割:" + t)
		}
		item := strings.ToLower(kv[0])
		if item == "@param" {
			for i := int64(0); ; i++ {
				item = "@param" + strconv.FormatInt(i, 10)
				if _, ok := commentMap[item]; !ok {
					commentMap[item] = kv[1]
					break
				}
			}
		} else if _, ok := commentMap[item]; !ok {
			commentMap[item] = kv[1]
		} else {
			return nil, errors.New("存在重复项 " + item)
		}
	}
	return commentToCode(commentMap)
}

func commentToCode(comments map[string]string) (map[string]interface{}, error) {
	// fmt.Println(comments)
	data := make(map[string]interface{})

	//如果无namespace则当做不处理
	if _, ok := comments["@namespace"]; !ok {
		return nil, nil
	}
	if _, ok := comments["@title"]; !ok {
		return nil, errors.New("@title不能为空")
	}
	if _, ok := comments["@return"]; !ok {
		return nil, errors.New("@return不能为空")
	}

	inputs := []map[string]interface{}{}
	for i := int64(0); ; i++ {
		item := "@param" + strconv.FormatInt(i, 10)
		if p, ok := comments[item]; !ok {
			break
		} else {
			items := strings.Fields(p)
			if len(items) <= 1 {
				return nil, errors.New("@param 描述信息必须有2项，必须依次是： name type [required] [desc]")
			}
			paramInfo := make(map[string]interface{})
			paramInfo["name"] = items[0]
			paramInfo["typ"] = strings.Title(items[1])
			paramInfo["required"] = false
			paramInfo["desc"] = ""
			if len(items) > 2 {
				if items[2] == "required" {
					paramInfo["required"] = true
				}
				paramInfo["desc"] = strings.TrimLeft(strings.Join(items[2:], ""), "required")
			}
			inputs = append(inputs, paramInfo)
		}

	}
	data["inputs"] = inputs
	for k, v := range comments {
		data[strings.TrimLeft(k, "@")] = v
	}
	return data, nil
}
