package ecs

import (
	"context"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepMountAlicloudDisk struct {
}

func (s *stepMountAlicloudDisk) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*ClientWrapper)
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)
	instance := state.Get("instance").(*ecs.Instance)

	alicloudDiskDevices := config.ECSImagesDiskMappings
	if len(config.ECSImagesDiskMappings) == 0 {
		return multistep.ActionContinue
	}

	ui.Say("Mounting disks.")

	describeDisksRequest := ecs.CreateDescribeDisksRequest()
	describeDisksRequest.RegionId = instance.RegionId
	describeDisksRequest.InstanceId = instance.InstanceId
	diskResponse, err := client.DescribeDisks(describeDisksRequest)
	if err != nil {
		return halt(state, err, "Error querying disks")
	}

	disks := diskResponse.Disks.Disk
	for _, disk := range disks {
		if disk.Status == DiskStatusAvailable {
			attachDiskRequest := ecs.CreateAttachDiskRequest()
			attachDiskRequest.DiskId = disk.DiskId
			attachDiskRequest.InstanceId = instance.InstanceId
			attachDiskRequest.Device = getDevice(&disk, alicloudDiskDevices)
			if _, err := client.AttachDisk(attachDiskRequest); err != nil {
				return halt(state, err, "Error mounting disks")
			}
		}
	}

	for _, disk := range disks {
		_, err = client.WaitForExpected(&WaitForExpectArgs{
			RequestFunc: func() (responses.AcsResponse, error) {
				request := ecs.CreateDescribeDisksRequest()
				request.RegionId = instance.RegionId
				request.DiskIds = disk.DiskId
				return client.DescribeDisks(request)
			},
			EvalFunc: func(response responses.AcsResponse, err error) WaitForExpectEvalResult {
				if err != nil {
					return WaitForExpectToRetry
				}

				disksResponse := response.(*ecs.DescribeDisksResponse)
				disks := disksResponse.Disks.Disk
				for _, disk := range disks {
					if disk.Status == DiskStatusInUse {
						return WaitForExpectSuccess
					}
				}

				return WaitForExpectToRetry
			},
			RetryTimes: defaultRetryTimes,
		})

		if err != nil {
			return halt(state, err, "Timeout waiting for mount")
		}
	}

	ui.Say("Finished mounting disks.")
	return multistep.ActionContinue
}

func (s *stepMountAlicloudDisk) Cleanup(state multistep.StateBag) {

}

func getDevice(disk *ecs.Disk, diskDevices []AlicloudDiskDevice) string {
	if disk.Device != "" {
		return disk.Device
	}
	for _, alicloudDiskDevice := range diskDevices {
		if alicloudDiskDevice.DiskName == disk.DiskName || alicloudDiskDevice.SnapshotId == disk.SourceSnapshotId {
			return alicloudDiskDevice.Device
		}
	}
	return ""
}
