package main

import (
	"fmt"
	"log"
	"flag"
	"github.com/wallnerryan/fli-docker/utils"
)

func main() {
    var user string
    var token string
    var endpoint string
    var manifest string
    var composeOpts string 

    // `docker-compose` path
    // TODO configurable in the future
    var composePath string
    composePath = "/usr/local/bin/docker-compose"

    // `dpcli` Path
    // TODO maybe configurable in the future
    var dpcliPath string
    dpcliPath = "/opt/clusterhq/bin/dpcli"

    // Check if needed dependencies are available
    isComposeAvail, err := utils.CheckForTool(composePath)
    if (!isComposeAvail){
    	fmt.Printf("-----------------------------------------------------------------------\n")
    	fmt.Printf("docker-compose is not installed, it is needed to use flitodock\n")
	fmt.Printf("docker-compose is available at https://docs.docker.com/compose/install/\n")
	fmt.Printf("-----------------------------------------------------------------------\n")
	log.Fatal(err.Error())
    }else{
	log.Println("docker-compose Ready!\n")
    }

    isDpcliAvail, err := utils.CheckForTool(dpcliPath)
    if (!isDpcliAvail){
    	fmt.Printf("-------------------------------------------------------\n")
    	fmt.Printf("dpcli is not installed, it is needed to use flitodock\n")
	fmt.Printf("dpcli is available at https://clusterhq.com\n")
	fmt.Printf("-------------------------------------------------------\n")
	log.Fatal(err.Error())
    }else{
	log.Println("dpcli Ready!\n")
    }

    flag.StringVar(&user, "u", "", "Flocker Hub username")
    flag.StringVar(&token, "t", "", "Flocker Hub user token")
    flag.StringVar(&endpoint, "v", "", "Flocker Hub endpoint")
    flag.StringVar(&manifest, "f", "manifest.yml", "Stateful application manifest file")
    /* 
    Im thinking this should be optional meaning if its not
    present then flidock will not also run the docker-compose command
    but rather will just edit the docker-compose.yml file in place
    and let the use run the docker-compose command. 
    This may be even a good option to start with instead of using
    '-c' at all.
    */
    flag.StringVar(&composeOpts, "c", "up", "Options to pass to Docker Compose such as 'up -d'")

    flag.Parse()

    /*
    # only for debug
    fmt.Printf("user = %s\n", user)
    fmt.Printf("token = %s\n", token)
    fmt.Printf("endpoint = %s\n", endpoint)
    fmt.Printf("manifest = %s\n", manifest)
    fmt.Printf("composeOpts = %s\n", composeOpts)
    */
}
