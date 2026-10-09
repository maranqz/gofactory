package other

import "factory/declaredFactoriesOnlyWithFactory/nested"

//gofactory:factory
func Build() nested.WithDeclaredFactory { // want Build:"gofactory:factory"
	return nested.WithDeclaredFactory{}
}
