package ecs

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepRegionCopyAlicloudImage struct {
	AlicloudImageDestinationRegions []string
	AlicloudImageDestinationNames   []string
	RegionId                        string
	ImageCopyEncrypted              bool
}

func (s *stepRegionCopyAlicloudImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if len(s.AlicloudImageDestinationRegions) == 0 && s.ImageCopyEncrypted == false {
		return multistep.ActionContinue
	} else {
		s.AlicloudImageDestinationRegions = append(s.AlicloudImageDestinationRegions, s.RegionId)
	}

	client := state.Get("client").(*ClientWrapper)
	imageId := state.Get("alicloudimage").(string)
	alicloudImages := state.Get("alicloudimages").(map[string]string)
	ui := state.Get("ui").(packer.Ui)

	numberOfName := len(s.AlicloudImageDestinationNames)

	ui.Say("Coping image...")
	for index, destinationRegion := range s.AlicloudImageDestinationRegions {

		ecsImageName := ""
		if numberOfName > 0 && index < numberOfName {
			ecsImageName = s.AlicloudImageDestinationNames[index]
		}

		copyImageRequest := ecs.CreateCopyImageRequest()
		copyImageRequest.RegionId = s.RegionId
		copyImageRequest.ImageId = imageId
		copyImageRequest.Encrypted = requests.Boolean(strconv.FormatBool(s.ImageCopyEncrypted))
		copyImageRequest.DestinationRegionId = destinationRegion
		copyImageRequest.DestinationImageName = ecsImageName

		image, err := client.CopyImage(copyImageRequest)
		if err != nil {
			return halt(state, err, "Error copying images")
		}

		if s.ImageCopyEncrypted == true {
			if _, err := client.WaitForImageStatus(destinationRegion, image.ImageId, ImageStatusAvailable, time.Duration(ALICLOUD_DEFAULT_LONG_TIMEOUT)*time.Second); err != nil {
				return halt(state, err, "error waiting copy images")
			}
		}

		alicloudImages[destinationRegion] = image.ImageId
		ui.Message(fmt.Sprintf("Copy image from %s(%s) to %s(%s)", s.RegionId, imageId, destinationRegion, image.ImageId))
	}
	return multistep.ActionContinue
}

func (s *stepRegionCopyAlicloudImage) Cleanup(state multistep.StateBag) {
	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if cancelled || halted {
		ui := state.Get("ui").(packer.Ui)
		client := state.Get("client").(*ClientWrapper)
		alicloudImages := state.Get("alicloudimages").(map[string]string)
		ui.Say(fmt.Sprintf("Stopping copy image because cancellation or error..."))
		for copiedRegionId, copiedImageId := range alicloudImages {
			if copiedRegionId == s.RegionId {
				continue
			}

			cancelCopyImageRequest := ecs.CreateCancelCopyImageRequest()
			cancelCopyImageRequest.RegionId = copiedRegionId
			cancelCopyImageRequest.ImageId = copiedImageId
			if _, err := client.CancelCopyImage(cancelCopyImageRequest); err != nil {
				ui.Say(fmt.Sprintf("Error cancelling copy image: %v", err))
			}
		}
	}
}
