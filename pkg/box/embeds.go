package box

//go:generate go run boxgen.go -compress templates templates/
//go:generate go run boxgen.go -compress static static/
//go:generate go run boxgen.go -compress dist dist/
//go:generate go run boxgen.go -compress -constraints "box boxconfig" config simpleauth.default.yml
