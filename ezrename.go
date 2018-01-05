package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattn/anko/vm"
	"github.com/metakeule/fmtdate"
	"github.com/spf13/cobra"
	"github.com/valyala/fasttemplate"
)

func addFunc(env *vm.Env) {
	env.Define("upper", strings.ToUpper)
	env.Define("lower", strings.ToLower)
	env.Define("repeat", strings.Repeat)
	env.Define("replace", strings.Replace)
	env.Define("trim", strings.Trim)
	env.Define("date", func(fm string) string {
		return fmtdate.Format(fm, time.Now())
	})
}

func doRename(filter, formula string, test bool) error {
	env := vm.NewEnv()
	addFunc(env)

	wd, _ := os.Getwd()

	return filepath.Walk(wd, func(fpath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		env.Define("name", f.Name())
		env.Define("full", fpath)
		env.Define("base", filepath.Base(f.Name()))
		env.Define("ext", filepath.Ext(f.Name()))
		env.Define("dir", filepath.Base(filepath.Dir(fpath)))

		t, err := fasttemplate.NewTemplate(formula, "{", "}")
		if err != nil {
			log.Fatalf("syntax error %s", err)
		}

		s := t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
			val, err := env.Execute(tag)
			if err != nil {
				return w.Write([]byte(tag))
			}
			return w.Write([]byte(fmt.Sprintf("%v", val)))
		})

		fmt.Println(fmt.Sprintf("%-40s\t=> %s", f.Name(), s))
		return nil
	})
}

func main() {

	testMode := false
	filter := "*"

	var cmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Usage()
				return
			}

			formula := strings.Join(args, " ")
			doRename(filter, formula, testMode)
		},
	}

	cmd.Flags().BoolVarP(&testMode, "test", "t", false, "just test,not run.")
	cmd.Flags().StringVarP(&filter, "filter", "f", "*", "files filter.")
	cmd.Execute()
}
