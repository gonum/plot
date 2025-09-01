module gonum.org/v1/plot

go 1.24.0

require (
	codeberg.org/go-fonts/latin-modern v0.4.0
	codeberg.org/go-fonts/liberation v0.5.0
	codeberg.org/go-latex/latex v0.2.0
	codeberg.org/go-pdf/fpdf v0.11.1
	git.sr.ht/~sbinet/gg v0.7.0
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b
	golang.org/x/image v0.30.0
	gonum.org/v1/gonum v0.16.0
	rsc.io/pdf v0.1.1
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	gonum.org/v1/tools v0.0.0-20200318103217-c168b003ce8c // indirect
)

tool (
	golang.org/x/tools/cmd/goimports
	gonum.org/v1/tools/cmd/check-copyright
	gonum.org/v1/tools/cmd/check-imports
)
