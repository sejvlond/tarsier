package persistent_storage

import (
	"fmt"
	"os"
	"sync"

	"github.com/sejvlond/tarsier/human_size"
	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewRead(lgr tarsier_plugins.Logger) *Read {
	return &Read{
		lgr: lgr,
	}
}

type Read struct {
	lgr tarsier_plugins.Logger
}

type ReadData struct {
	File       string
	Amount     string
	Concurrent uint
}

func (cmd *Read) Name() string {
	return NAME + "/read"
}

func (cmd *Read) Description() string {
	return `Read will read as many data as you say from specify file in paralel<br>
Data:
<pre>
	amount: how much data to read (number in bytes, or string with GB, MB, kB, b)
	file: file from data will be read
	concurrent: number of paralel readers
</pre>`
}

func (cmd *Read) DataStruct() interface{} {
	return &ReadData{}
}

func (cmd *Read) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*ReadData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	amount, ok := human_size.Parse(data.Amount)
	if !ok {
		return "Invalid amount", fmt.Errorf("Invalid amount")
	}
	var wg sync.WaitGroup
	var glErr error
	for i := uint(0); i < data.Concurrent; i++ {
		wg.Add(1)
		go func() {
			cmd.lgr.Infof("Starting to read")
			file, err := os.Open(data.File)
			if err != nil {
				glErr = err
			} else {
				buf := make([]byte, 1)
				for i := 0; i < amount; i++ {
					_, err := file.Read(buf)
					if err != nil {
						glErr = err
						break
					}
				}
			}
			file.Close()
			wg.Done()
			cmd.lgr.Infof("Everything was read")
		}()
	}
	wg.Wait()
	if glErr != nil {
		return fmt.Sprintf("Error reading from file: %v", glErr), glErr
	}
	return fmt.Sprintf("%v of text was read from %v file in %v paralel",
		data.Amount, data.File, data.Concurrent), nil
}
