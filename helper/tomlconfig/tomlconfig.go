package tomlconfig

import (
		"log"
		"io/ioutil"
		"os"
		"io"
)

type Parser interface {
		ParseTOML([]byte) error
}

func Load(configFile string, p Parser) {
		var err error
		var input = io.ReadCloser(os.Stdin)
		if input, err = os.Open(configFile); err != nil {
				log.Fatalln(err)
		}

		// Read the config file
		tomlBytes, err := ioutil.ReadAll(input)
		input.Close()
		if err != nil {
				log.Fatalln(err)
		}

		// Parse the config
		if err := p.ParseTOML(tomlBytes); err != nil {
				log.Fatalln("Could not parse %q: %v", configFile, err)
		}
}
