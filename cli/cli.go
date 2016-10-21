package cli

import (
	"os/exec"
	"fmt"

	"github.com/ClusterHQ/fli-docker/types"
	"github.com/ClusterHQ/fli-docker/logger"
	"github.com/ClusterHQ/fli-docker/utils"
)

/*
	Bindings to the FlockerHub CLI
*/

func SetFlockerHubEndpoint(endpoint string, bin bool) {
	var fli string
	if (bin){
		fli = "fli "
	}else{
		fli = utils.FliDockerCmd
	}

	logger.Info.Println("Setting FlockerHub Endpoint: ", endpoint)
	var cmd = fmt.Sprintf("%s config -u %s", fli, endpoint)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not set endpoint")
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

func GetFlockerHubEndpoint(bin bool) (flockerhubEndpoint string, err error) {
	var fli string
	if (bin){
		fli = "fli "
	}else{
		fli = utils.FliDockerCmd
	}

	logger.Info.Println("Getting FlockerHub Endpoint")
	var cmd = fmt.Sprintf("%s config | grep 'FlockerHub URL:' | awk '{print $3}'", fli)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not get endpoint")
		logger.Error.Println(err)
		return "", err
	}
	logger.Info.Println(string(out))
	return string(out), nil
}

func SetFlockerHubTokenFile(tokenFile string, bin bool) {
	var fli string
	if (bin){
		fli = "fli "
	}else{
		fli = utils.FliDockerCmd
	}

	logger.Info.Println("Setting FlockerHub Tokenfile: ", tokenFile)
	var cmd = fmt.Sprintf("%s config -t %s", fli, tokenFile)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not set tokenfile")
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

func GetFlockerHubTokenFile(bin bool) (flockerHubTokenFile string, err error) {
	var fli string
	if (bin){
		fli = "fli "
	}else{
		fli = utils.FliDockerCmd
	}

	logger.Info.Println("Getting FlockerHub Tokenfile")
	var cmd = fmt.Sprintf("%s config | grep 'Authentication Token File:' | awk '{print $3}'", fli)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not get tokenfile")
		logger.Error.Println(err)
		return "", err
	}
	logger.Info.Println(string(out))
	return string(out), nil
}

// Run the command to sync a volumeset
func syncVolumeset(volumeSetId string, bin bool) {
	var fli string
	if (bin){
		fli = "fli "
	}else{
		fli = utils.FliDockerCmd
	}

	logger.Info.Println("Syncing Volumeset: ", volumeSetId)
	var cmd = fmt.Sprintf("%s sync %s", fli, volumeSetId)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not sync dataset")
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

// Run the command to pull a specific snapshot
func pullSnapshot(volumeSetId string, snapshotId string, bin bool){
	var fli string
	if (bin){
		fli = "fli "
	}else{
		fli = utils.FliDockerCmd
	}

	logger.Info.Println("Pulling Snapshot: ", snapshotId)
	var cmd = fmt.Sprintf("%s pull %s:%s", fli, volumeSetId, snapshotId)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not pull dataset, reason")
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

// Wrapper for sync and pull which takes
// a list of type Volume
func PullSnapshots(volumes []types.Volume, bin bool) {
	for _, volume := range volumes {
		syncVolumeset(volume.VolumeSet, bin)
		pullSnapshot(volume.VolumeSet, volume.Snapshot, bin)
	}
}

// Created a volume and returns it.
func createVolumeFromSnapshot(volumeName string, volumeSet string, snapshotId string, bin bool) (vol types.NewVolume, err error){
	var fli string
	if (bin){
		fli = "fli "
	}else{
		fli = utils.FliDockerCmd
	}

	logger.Info.Println("Creating Volume from Snapshot: ", snapshotId)
	var attrString = fmt.Sprintf("created_by=fli-docker,from_snap=%s", snapshotId)
	uuid, err := utils.GenUUID()
	if err != nil {
		logger.Error.Fatal(err)
	}

	var volName = fmt.Sprintf("fli-%s", uuid)
	var createCmd = fmt.Sprintf("%s clone %s:%s -a %s %s", fli, volumeSet, snapshotId, attrString, volName)
	cmd := exec.Command("sh", "-c", createCmd)
	createOut, err := cmd.Output()
	if err != nil {
		logger.Error.Fatal(err)
	}
	var path = string(createOut)
	logger.Info.Println(path)
	if path == "" {
			logger.Error.Fatal("Could not find volume path")
	 }
	return types.NewVolume{Name: volumeName, VolumePath: path}, nil
}

func CreateVolumesFromSnapshots(volumes []types.Volume, bin bool) (newVols []types.NewVolume, err error) {
	vols := []types.NewVolume{}
	for _, volume := range volumes {
		vol, err := createVolumeFromSnapshot(volume.Name, volume.VolumeSet, volume.Snapshot, bin)
		if err != nil {
			return nil, err
		}else {
			vols = append(vols, vol)
		}
	}
	return vols, nil
}