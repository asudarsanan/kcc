package list

import (
	"fmt"
	"kcc-cli/m/v2/resources"
)

func Contexts(contexts []resources.Context) {
	for _, context := range contexts {
		fmt.Println(context.Name)
	}
}
