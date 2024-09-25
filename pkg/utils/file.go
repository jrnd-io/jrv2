package utils

import "os"

func CountFilesInDir(dir string) (int, error) {

	existentDir, _ := Exists(dir)
	if !existentDir {
		return 0, nil
	}

	f, err := os.Open(dir)
	if err != nil {
		return 0, err
	}
	list, err := f.Readdirnames(-1)
	if err != nil {
		return 0, err
	}

	err = f.Close()
	if err != nil {
		return 0, err
	}
	return len(list), nil
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
