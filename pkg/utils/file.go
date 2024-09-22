package utils

import "os"

func CountFilesInDir(dir string) int {

	existentDir, _ := Exists(dir)
	if !existentDir {
		return 0
	}

	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	list, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
	return len(list)
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
