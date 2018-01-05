# ezrename
`ezrename` is a command-line rename tool written in [Go](https://golang.org/). `ezrename` using [anko](https://github.com/mattn/anko) as script. Yes, you can coding in your filename.

## UseAge
    ezrename --path=./ {script}some_const_text{script}{script}

## example
    ezrename --path=test {date('YYMM')}/{index+1}/{base[0:2]}_{randstr(3)}{ext} 
    test\somedir                         => test\1801\1\so_Itc
    test\00_1.txt                        => test\1801\2\00_WZ4.txt
    test\00_2.txt                        => test\1801\3\00_26z.txt
    test\01_3.txt                        => test\1801\4\01_Aq1.txt

## functions
    upper   func(string)string
    lower   func(string)string
    repeat  func(string)string
    replace func(string)string
    trim    func(string)string
    date    func(string)string
    format  func(string, ...)string
    rand    func()int
    randstr func(int)string

## vars
    name    a/b/cc.txt -> cc.txt
    full    a/b/cc.txt -> a/b/cc.txt
    base    a/b/cc.txt -> cc
    ext     a/b/cc.txt -> .txt
    dir     a/b/cc.txt -> b
    index   index number 
