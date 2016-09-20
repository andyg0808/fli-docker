package main

import (
    "fmt"
    "log"
    "flag"
    "path/filepath"
    "io/ioutil"
    "github.com/wallnerryan/fli-docker/utils"
)

func main() {

    // should this be a struct?
    var user string
    var token string
    var endpoint string
    var manifest string
    var composeOpts string 

    var composeCmd string
    composeCmd = "docker-compose version"

    var fliCmd string
    fliCmd = "dpcli init" //this will need `fli version` or somthing

    // Check if needed dependencies are available
    isComposeAvail, err := utils.CheckForCmd(composeCmd)
    if (!isComposeAvail){
        fmt.Printf("-----------------------------------------------------------------------\n")
        fmt.Printf("docker-compose is not installed, it is needed to use fli-docker\n")
        fmt.Printf("docker-compose is available at https://docs.docker.com/compose/install/\n")
        fmt.Printf("-----------------------------------------------------------------------\n")
        log.Fatal(err.Error())
    }else{
        log.Println("docker-compose Ready!\n")
    }

    isFliAvail, err := utils.CheckForCmd(fliCmd)
    if (!isFliAvail){
        fmt.Printf("-------------------------------------------------------\n")
        fmt.Printf("fli is not installed, it is needed to use fli-docker\n")
        fmt.Printf("fli is available at https://clusterhq.com\n")
        fmt.Printf("-------------------------------------------------------\n")
    }else{
        log.Println("fli Ready!\n")
    }

    flag.StringVar(&user, "u", "", "Flocker Hub username")
    flag.StringVar(&token, "t", "", "Flocker Hub user token")
    // Should we replace or add the above with the option to point to vhub.txt?
    flag.StringVar(&endpoint, "v", "", "Flocker Hub endpoint")
    flag.StringVar(&manifest, "f", "manifest.yml", "Stateful application manifest file")
    flag.StringVar(&composeOpts, "c", "up", "Options to pass to Docker Compose such as 'up -d'") //optional

    // Parse all the flags from user input
    flag.Parse()

    /*
    # only for debug
    fmt.Printf("user = %s\n", user)
    fmt.Printf("token = %s\n", token)
    fmt.Printf("endpoint = %s\n", endpoint)
    fmt.Printf("manifest = %s\n", manifest)
    fmt.Printf("composeOpts = %s\n", composeOpts)
    */

    //TODO check for empty vars, or default ones.

    // 1. Verify that the manifest exists
    isManifestAvail, err := utils.CheckForFile(manifest)
    if (!isManifestAvail){
        log.Fatal(err.Error())
    }else{
        log.Println("fli-docker manifest found\n")
    }

    // 2. Process the manifest into a Struct in YAML
    //    and get a mapping of everything including:
    //         comopse_file: <file> (from `docker_app`)
    //         compose_volume_name : {volumeset: <id>, snapshot: <id>} (from `volumes`)
    //         flocker_hub : {endpoint : <url>, auth_token : <token>} (from `flocker_hub`), token can be optional
    
    // Get the yaml file passed in the args.
    filename, _ := filepath.Abs(manifest)
    // Read the file.
    yamlFile, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    // Pass the file to the ParseManifest
    fmt.Printf("Trying to unmarshall yaml file\n")
    m := utils.ParseManifest(yamlFile)

    fmt.Println(m.DockerApp)
    fmt.Println(m.Hub)
    fmt.Println(m.Hub.Endpoint)
    fmt.Println(m.Hub.AuthToken)
    fmt.Println(m.Volumes)
    fmt.Println(m.Volumes[0])
    fmt.Println(m.Volumes[1])

    // 3. Verify that the compose file exists.
    isComposeFileAvail, err := utils.CheckForFile(m.DockerApp)
    if (!isComposeFileAvail){
        log.Fatal(err.Error())
    }else{
        log.Println("docker-compose file found\n")
    }

    // 4. Try and pull snapshots
    // TODO need to return err and check for it
    utils.PullSnapshots(m.Volumes)

    // 5. Create volumes from snapshots and map them to 
    //    newVolPaths = {compose_volume_name : "/chq/<vol_path>"}
    newVolPaths, err := utils.CreateVolumesFromSnapshots(m.Volumes)
    fmt.Println(newVolPaths)

    // TODO need to return err and check for it
    utils.MakeCopy(m.DockerApp)

    // 7. replace volume_name with volume_name's associated "/chq/<vol_path/"
    // 8. write file back to compose file

    // TODO Loop through newVolPaths and perform
    //utils.MapVolumeToCompose(<volname>, <chq_vol_path>, m.DockerApp)

    // Verify Compose File
    utils.ParseCompose(m.DockerApp)
    // 9. (IF) -c is there for compose args, run compose, if not, done.

}