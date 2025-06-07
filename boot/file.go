package boot

import "os"

func Initfile() error {
	if _, err := os.Stat("judges"); os.IsNotExist(err) {
		if err := os.Mkdir("judges", 0755); err != nil {
			return err
		}
	}
	return nil
}
