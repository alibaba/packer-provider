package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/mitchellh/packer/template/interpolate"
)

type AlicloudDiskDevice struct {
	DiskName           string `mapstructure:"alicloud_disk_name"`
	DiskCategory       string `mapstructure:"alicloud_disk_category"`
	DiskSize           int    `mapstructure:"alicloud_disk_size"`
	SnapshotId         string `mapstructure:"alicloud_disk_snapshot_id"`
	Description        string `mapstructure:"alicloud_disk_description"`
	DeleteWithInstance bool   `mapstructure:"alicloud_disk_delete_with_instance"`
	Device             string `mapstructure:"alicloud_disk_device"`
}

type AlicloudDiskDevices struct {
	ECSImagesDiskMappings []AlicloudDiskDevice `mapstructure:"alicloud_image_disk_mappings"`
}

type AlicloudImageConfig struct {
	AlicloudImageName                 string   `mapstructure:"alicloud_image_name"`
	AlicloudImageVersion              string   `mapstructure:"alicloud_image_version"`
	AlicloudImageDescription          string   `mapstructure:"alicloud_image_description"`
	AlicloudImageShareAccounts        []string `mapstructure:"alicloud_image_share_account"`
	AlicloudImageUNShareAccounts      []string `mapstructure:"alicloud_image_unshare_account"`
	AlicloudImageDestinationRegions   []string `mapstructure:"alicloud_image_copy_regions"`
	AlicloudImageDestinationNames     []string `mapstructure:"alicloud_image_copy_names"`
	AlicloudImageForceDetele          bool     `mapstructure:"alicloud_image_force_delete"`
	AlicloudImageForceDeteleSnapshots bool     `mapstructure:"alicloud_image_force_delete_snapshots"`
	AlicloudImageForceDeleteInstances bool     `mapstructure:"alicloud_image_force_delete_instances"`
	AlicloudImageSkipRegionValidation bool     `mapstructure:"alicloud_skip_region_validation"`
	AlicloudDiskDevices               `mapstructure:",squash"`
}

func (c *AlicloudImageConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error
	if c.AlicloudImageName == "" {
		errs = append(errs, fmt.Errorf("alicloud_image_name must be specified"))
	}

	if len(c.AlicloudImageDestinationRegions) > 0 {
		regionSet := make(map[string]struct{})
		regions := make([]string, 0, len(c.AlicloudImageDestinationRegions))

		for _, region := range c.AlicloudImageDestinationRegions {
			// If we already saw the region, then don't look again
			if _, ok := regionSet[region]; ok {
				continue
			}

			// Mark that we saw the region
			regionSet[region] = struct{}{}

			if !c.AlicloudImageSkipRegionValidation {
				// Verify the region is real
				if valid := validateRegion(region); valid != nil {
					errs = append(errs, fmt.Errorf("Unknown region: %s", region))
					continue
				}
			}

			regions = append(regions, region)
		}

		c.AlicloudImageDestinationRegions = regions
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateRegion(region string) error {

	for _, valid := range common.ValidRegions {
		if region == string(valid) {
			return nil
		}
	}

	return fmt.Errorf("Not a valid alicloud region: %s", region)
}
