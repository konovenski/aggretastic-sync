// Command aggretastic-generator is used for automated managing
// https://github.com/AaHaInc/Aggretastic repository.
//
// Aggretastic-generator can update currently used,
// add new and remove deprecated aggregations
//
// All source-related operations are performed via
// high-level go/ast wrapper: pretty_dst.
//
// Aggretastic-generator can be used for tracking
// several upstreams for several Aggretastic versions.
// You can specify upstream version in conf.env file
// in your Aggretastic repo
//
package main

import (
	"github.com/joho/godotenv"
	"github.com/dkonovenschi/aggretastic-sync/olivere_v6_pipelines"
	"log"
	"os"
)

func main() {
	err := godotenv.Load("conf.env")
	if err != nil {
		err = godotenv.Load("conf.env.default")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	switch os.Getenv("REPO_BRANCH") {
	case "refs/heads/release-branch.v6":
		olivere_v6_pipelines.Run()
	default:
		panic("Can't find any pipelines for this branch")
	}

}
