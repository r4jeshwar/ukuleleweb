//go:build linux

package main

import (
	"log"

	"github.com/landlock-lsm/go-landlock/landlock"
)

func restrictAccess(rwDirs ...string) {
	err := landlock.V2.BestEffort().RestrictPaths(
		landlock.RWDirs(rwDirs...),
	)
	if err != nil {
		log.Fatalf("Landlock: %v", err)
	}
}
