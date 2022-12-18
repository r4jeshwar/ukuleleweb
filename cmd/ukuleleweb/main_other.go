//go:build !linux

package main

func restrictAccess(rwDirs ...string) {
	// Noop
}
