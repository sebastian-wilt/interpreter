package interpret

import "fmt"

// Inbuilt
type Inbuilt struct {
	name string
}

func (i *Inbuilt) Name() string {
	return i.name
}

func (i *Inbuilt) Type() {}

func getInbuilts() map[string]Type {
	inbuilts := map[string]Type{}
	types := []string{"int", "string", "real", "char", "boolean"}
	for _, s := range types {
		inbuilts[s] = &Inbuilt{name: s}
	}

	return inbuilts
}

func getInbuiltValue(i *Inbuilt) Value {
	switch i.name {
	case "int":
		return &Integer{Value: 0}
	case "real":
		return &Real{Value: 0.0}
	case "string":
		return &String{Value: ""}
	case "char":
		return &Char{Value: '\000'}
	case "boolean":
		return &Boolean{Value: false}
	default:
		panic(fmt.Sprintf("Unknown inbuilt: %s\n", i.name))
	}
}
