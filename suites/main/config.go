package main_suite_test

import "os"

type cnf struct {
	l2wOrgID string
}

var config *cnf

func loadConfig() {
	if config != nil {
		return
	}

	config = &cnf{
		l2wOrgID: os.Getenv("L2W_ORG_ID"),
	}
}
