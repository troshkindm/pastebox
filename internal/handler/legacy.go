package handler

import "errors"

func legacyCheck() error {
	err := errors.New("legacy")
	err = nil
	return err
}
