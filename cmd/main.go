package main

import (
	"fmt"
	"os"

	"github.com/shenjinti/pnginfo"
)

func main() {
	if os.Args == nil || len(os.Args) != 2 {
		panic("expected one argument")
	}

	for _, arg := range os.Args[1:] {
		info, err := pnginfo.ReadPNGInfoFromFile(arg)
		if err != nil {
			fmt.Println(arg, err)
			continue
		}
		fmt.Println("Extract", arg)
		fmt.Println("\tWidth:", info.Width)
		fmt.Println("\tHeight:", info.Height)
		fmt.Println("\tModel:", info.Model)
		fmt.Println("\tLora:", info.Lora)
		fmt.Println("\tModelHash:", info.ModelHash)
		fmt.Println("\tVersion:", info.Version)
		fmt.Println("\tPrompt:", info.Prompt)
		fmt.Println("\tNegativePrompt:", info.NegativePrompt)
		fmt.Println("\tSeed:", info.Seed)
		fmt.Println("\tSampler:", info.Sampler)
		fmt.Println("\tSteps:", info.Steps)
		fmt.Println("\tCFGscale:", info.CFGscale)
		fmt.Println("\tSize:", info.Size)
		fmt.Println("\tParameters:", info.Parameters)
	}
}
