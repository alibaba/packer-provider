package ecs

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/packer/packer"
)

type Artifact struct {
	// A map of regions to alicloud image IDs.
	AlicloudImages map[string]string

	// BuilderId is the unique ID for the builder that created this alicloud image
	BuilderIdValue string

	// Alcloud connection for performing API stuff.
	Client *ClientWrapper
}

func (a *Artifact) BuilderId() string {
	return a.BuilderIdValue
}

func (*Artifact) Files() []string {
	// We have no files
	return nil
}

func (a *Artifact) Id() string {
	parts := make([]string, 0, len(a.AlicloudImages))
	for region, ecsImageId := range a.AlicloudImages {
		parts = append(parts, fmt.Sprintf("%s:%s", region, ecsImageId))
	}

	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func (a *Artifact) String() string {
	alicloudImageStrings := make([]string, 0, len(a.AlicloudImages))
	for region, id := range a.AlicloudImages {
		single := fmt.Sprintf("%s: %s", region, id)
		alicloudImageStrings = append(alicloudImageStrings, single)
	}

	sort.Strings(alicloudImageStrings)
	return fmt.Sprintf("Alicloud images were created:\n\n%s", strings.Join(alicloudImageStrings, "\n"))
}

func (a *Artifact) State(name string) interface{} {
	switch name {
	case "atlas.artifact.metadata":
		return a.stateAtlasMetadata()
	default:
		return nil
	}
}

func (a *Artifact) Destroy() error {
	errors := make([]error, 0)

	for region, imageId := range a.AlicloudImages {
		log.Printf("Delete alicloud image ID (%s) from region (%s)", imageId, region)

		// Get alicloud image metadata
		describeImagesRequest := ecs.CreateDescribeImagesRequest()
		describeImagesRequest.RegionId = region
		describeImagesRequest.ImageId = imageId
		imagesResponse, err := a.Client.DescribeImages(describeImagesRequest)
		if err != nil {
			errors = append(errors, err)
		}

		images := imagesResponse.Images.Image
		if len(images) == 0 {
			err := fmt.Errorf("Error retrieving details for alicloud image(%s), no alicloud images found ", imageId)
			errors = append(errors, err)
			continue
		}

		//Unshared the shared account before destroy
		describeImageSharePermissionRequest := ecs.CreateDescribeImageSharePermissionRequest()
		describeImageSharePermissionRequest.RegionId = region
		describeImageSharePermissionRequest.ImageId = imageId
		imagesSharePermissionResponse, err := a.Client.DescribeImageSharePermission(describeImageSharePermissionRequest)
		if err != nil {
			errors = append(errors, err)
		}

		accountsNumber := len(imagesSharePermissionResponse.Accounts.Account)
		if accountsNumber > 0 {
			accounts := make([]string, accountsNumber)
			for index, account := range imagesSharePermissionResponse.Accounts.Account {
				accounts[index] = account.AliyunId
			}

			modifyImageSharePermissionReq := ecs.CreateModifyImageSharePermissionRequest()
			modifyImageSharePermissionReq.RegionId = region
			modifyImageSharePermissionReq.ImageId = imageId
			modifyImageSharePermissionReq.RemoveAccount = &accounts
			_, err := a.Client.ModifyImageSharePermission(modifyImageSharePermissionReq)
			if err != nil {
				errors = append(errors, err)
			}
		}

		// Delete alicloud images
		deleteImageRequest := ecs.CreateDeleteImageRequest()
		deleteImageRequest.ImageId = imageId
		deleteImageRequest.RegionId = region
		if _, err := a.Client.DeleteImage(deleteImageRequest); err != nil {
			errors = append(errors, err)
		}

		//Delete the snapshot of this images
		for _, diskDevices := range images[0].DiskDeviceMappings.DiskDeviceMapping {
			deleteSnapshotRequest := ecs.CreateDeleteSnapshotRequest()
			deleteSnapshotRequest.SnapshotId = diskDevices.SnapshotId
			_, err := a.Client.DeleteSnapshot(deleteSnapshotRequest)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		if len(errors) == 1 {
			return errors[0]
		} else {
			return &packer.MultiError{Errors: errors}
		}
	}

	return nil
}

func (a *Artifact) stateAtlasMetadata() interface{} {
	metadata := make(map[string]string)
	for region, imageId := range a.AlicloudImages {
		k := fmt.Sprintf("region.%s", region)
		metadata[k] = imageId
	}

	return metadata
}
