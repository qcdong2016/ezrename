package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mattn/anko/vm"
	"github.com/metakeule/fmtdate"
	"github.com/spf13/cobra"
	"github.com/valyala/fasttemplate"
)

func NewEnv() *vm.Env {
	env := vm.NewEnv()

	env.Define("upper", strings.ToUpper)
	env.Define("lower", strings.ToLower)
	env.Define("repeat", strings.Repeat)
	env.Define("replace", strings.Replace)
	env.Define("trim", strings.Trim)
	env.Define("date", func(fm string) string {
		return fmtdate.Format(fm, time.Now())
	})
	env.Define("format", fmt.Sprintf)

	rand.Seed(time.Now().UnixNano())
	env.Define("rand", func() int {
		return rand.Intn(int(^uint(0) >> 1))
	})

	var bytes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	env.Define("randstr", func(length int) string {
		result := make([]byte, length)
		for i := 0; i < length; i++ {
			j := rand.Intn(len(bytes))
			result[i] = bytes[j]
		}
		return string(result)
	})

	return env
}

func basename(filename string) string {
	i := strings.LastIndex(filename, ".")
	if i != -1 {
		return filename[0:i]
	}
	return filename
}

func setEnv(env *vm.Env, file *FileInfo) {
	env.Define("name", filepath.Base(file.Name))
	env.Define("full", file.FullName)
	env.Define("base", basename(file.Name))
	env.Define("ext", filepath.Ext(file.Name))
	env.Define("dir", filepath.Base(filepath.Dir(file.FullName)))
	env.Define("index", file.Index)
	env.Define("isdir", file.Info.IsDir())
}

type FileInfo struct {
	Info        os.FileInfo
	Name        string
	FullName    string
	NewName     string
	NewFullName string
	Index       int
}

var config = struct {
	Test       bool
	Filter     string
	TargetPath string
	Formula    string
	Sort       string
}{}

func ListDir(dirPth string) (files []*FileInfo, err error) {
	files = []*FileInfo{}

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return
	}

	for _, fi := range dir {
		fullname := filepath.Join(dirPth, fi.Name())
		files = append(files, &FileInfo{
			FullName: fullname,
			Info:     fi,
			Name:     fi.Name(),
		})
	}

	return
}

func do() {
	files_ori, _ := ListDir(config.TargetPath)

	env := NewEnv()

	files := []*FileInfo{}
	if config.Filter != "" {
		for _, file := range files_ori {
			setEnv(env, file)
			val, err := env.Execute(config.Filter)
			if err != nil {
				log.Fatalf("syntax error [filter]:%s", err)
			}
			switch val.(type) {
			case bool:
				if val.(bool) {
					files = append(files, file)
				}
			default:
				log.Fatalf("filter should return a bool")
			}
		}
	} else {
		files = files_ori
	}

	if config.Sort != "" {
		sort.Slice(files, func(a, b int) bool {
			setEnv(env, files[a])
			va, err := env.Execute(config.Sort)
			if err != nil {
				log.Fatalf("runing sort script error: [%s]%s", files[a].FullName, err)
			}

			setEnv(env, files[b])
			vb, err := env.Execute(config.Sort)
			if err != nil {
				log.Fatalf("runing sort script error: [%s]%s", files[b].FullName, err)
			}

			switch va.(type) {
			case string:
			default:
				log.Fatalf("sort script shout return a string: [%s]", files[a].FullName)
			}

			switch vb.(type) {
			case string:
			default:
				log.Fatalf("sort script shout return a string: [%s]", files[b].FullName)
			}

			return va.(string) < vb.(string)
		})
	}

	for index, file := range files {
		file.Index = index
		setEnv(env, file)

		t, err := fasttemplate.NewTemplate(config.Formula, "{", "}")
		if err != nil {
			log.Fatalf("syntax error %s", err)
		}

		s := t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
			val, err := env.Execute(tag)
			if err != nil {
				log.Fatalf("syntax error [%s]:%s", file.Name, err)
			}

			return w.Write([]byte(fmt.Sprintf("%v", val)))
		})

		file.NewName = s
		file.NewFullName = filepath.Join(filepath.Dir(file.FullName), file.NewName)

		fmt.Println(fmt.Sprintf("%-40s\t=> %s", file.FullName, file.NewFullName))
		if !config.Test {
			dir := filepath.Dir(file.NewFullName)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				log.Fatal(err)
			}
			if err := os.Rename(file.FullName, file.NewFullName); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {

	config.TargetPath, _ = os.Getwd()

	var cmd = &cobra.Command{
		Use: "{script}",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Usage()
				return
			}

			config.Formula = strings.Join(args, " ")
			do()
		},
	}

	cmd.Flags().BoolVarP(&config.Test, "test", "t", false, "just print, not run.")
	cmd.Flags().StringVarP(&config.Filter, "filter", "f", "", "files filter script. bool")
	cmd.Flags().StringVarP(&config.Sort, "sort", "s", "", "files sort script. string")
	cmd.Flags().StringVarP(&config.TargetPath, "path", "p", config.TargetPath, "target path.")
	cmd.Execute()
}
