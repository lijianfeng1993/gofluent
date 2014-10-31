package main

import (
	"fmt"
	"os"
)

type Input interface {
	new() interface{}
	start(ctx chan Context) error
	configure(f map[string]string) error
}

var input_plugins = make(map[string]Input)

func RegisterInput(name string, input Input) {
	if input == nil {
		panic("input: Register input is nil")
	}

	if _, ok := input_plugins[name]; ok {
		panic("input: Register called twice for input " + name)
	}

	input_plugins[name] = input
}

func NewInputs(ctx chan Context) {
	router.AddInChan(ctx)

	for _, input_config := range config.Inputs_config {
		f := input_config.(map[string]string)
		go func(f map[string]string) {
			intput_type, ok := f["type"]
			if !ok {
				fmt.Println("no type configured")
				os.Exit(-1)
			}

			input, ok := input_plugins[intput_type]
			if !ok {
				fmt.Println("unkown type ", intput_type)
				os.Exit(-1)
			}

			in := input.new()

			err := in.(Input).configure(f)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			err = in.(Input).start(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		}(f)
	}

	return
}
