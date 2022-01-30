module github.com/psilva261/6to5

go 1.16

replace (
	github.com/tdewolff/parse/v2 v2.5.27 => github.com/psilva261/parse/v2 v2.5.28-0.20220130150813-6734237d078c
)

require (
	github.com/jvatic/goja-babel v0.0.0-20211030111852-5873af3b41cc
	github.com/tdewolff/parse/v2 v2.5.27 // indirect
)
