# ezrename
`ezrename` is a command-line rename tool written in [Go](https://golang.org/). `ezrename` using [anko](https://github.com/mattn/anko) script as expression. That powerful and flexible.

## UseAge
    ezrename --path=./ --sort=script --filter=script  {script}some_const_text{script}{script}

## example
![preview][]


## functions
    upper   func(string)string
    lower   func(string)string
    repeat  func(s string, count int)string
    replace func(s, old, new string, n int)string
    trim    func(s string, cutset string)string
    format  func(format string, ...)string
    rand    func()int
    randstr func(length int)string
    date    func(format string)string
help of [date](https://github.com/metakeule/fmtdate)

## vars
    name    a/b/cc.txt -> cc.txt
    full    a/b/cc.txt -> a/b/cc.txt
    base    a/b/cc.txt -> cc
    ext     a/b/cc.txt -> .txt
    dir     a/b/cc.txt -> b
    index   index number 
    isdir

[preview]: preview.jpg
