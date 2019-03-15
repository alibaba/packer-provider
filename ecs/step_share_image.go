package ecs

import (
	"context"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepShareAlicloudImage struct {
	AlicloudImageShareAccounts   []string
	AlicloudImageUNShareAccounts []string
	RegionId                     string
}

func (s *stepShareAlicloudImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*ClientWrapper)
	alicloudImages := state.Get("alicloudimages").(map[string]string)

	for copiedRegion, copiedImageId := range alicloudImages {
		modifyImageSharePermissionReq := ecs.CreateModifyImageSharePermissionRequest()
		modifyImageSharePermissionReq.RegionId = copiedRegion
		modifyImageSharePermissionReq.ImageId = copiedImageId
		modifyImageSharePermissionReq.AddAccount = &s.AlicloudImageShareAccounts
		modifyImageSharePermissionReq.RemoveAccount = &s.AlicloudImageUNShareAccounts

		if _, err := client.ModifyImageSharePermission(modifyImageSharePermissionReq); err != nil {
			return halt(state, err, "Failed modifying image share permissions")
		}
	}
	return multistep.ActionContinue
}

func (s *stepShareAlicloudImage) Cleanup(state multistep.StateBag) {
	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)

	if !cancelled && !halted {
		return
	}

	ui := state.Get("ui").(packer.Ui)
	client := state.Get("client").(*ClientWrapper)
	alicloudImages := state.Get("alicloudimages").(map[string]string)

	ui.Say("Restoring image share permission because cancellations or error...")

	for copiedRegion, copiedImageId := range alicloudImages {
		modifyImageSharePermissionRequest := ecs.CreateModifyImageSharePermissionRequest()
		modifyImageSharePermissionRequest.RegionId = copiedRegion
		modifyImageSharePermissionRequest.ImageId = copiedImageId
		modifyImageSharePermissionRequest.AddAccount = &s.AlicloudImageUNShareAccounts
		modifyImageSharePermissionRequest.RemoveAccount = &s.AlicloudImageShareAccounts
		if _, err := client.ModifyImageSharePermission(modifyImageSharePermissionRequest); err != nil {
			ui.Say(fmt.Sprintf("Restoring image share permission failed: %s", err))
		}
	}
}
