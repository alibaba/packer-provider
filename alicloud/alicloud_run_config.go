package alicloud

import (
	"errors"
	"fmt"
	"github.com/mitchellh/packer/common/uuid"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/template/interpolate"
	"os"
	"time"
)

const (
	ssh_time_out      = 60000000000
	default_port      = 22
	default_comm_type = "ssh"
)

type RunConfig struct {
	AssociatePublicIpAddress bool          `mapstructure:"alicloud_associate_public_ip_address"`
	ZoneId                   string        `mapstructure:"alicloud_zone_id"`
	IOOptimized              bool          `mapstructure:"alicloud_io_optimized"`
	InstanceType             string        `mapstructure:"alicloud_instance_type"`
	Description              string        `mapstructure:"alicloud_description"`
	AlicloudSourceImage      string        `mapstructure:"alicloud_source_image"`
	ForceStopInstance        bool          `mapstructure:"alicloud_force_stop_instance"`
	SecurityGroupId          string        `mapstructure:"alicloud_security_group_id"`
	SecurityGroupName        string        `mapstructure:"alicloud_security_group_name"`
	UserData                 string        `mapstructure:"alicloud_user_data"`
	UserDataFile             string        `mapstructure:"alicloud_user_data_file"`
	VpcId                    string        `mapstructure:"alicloud_vpc_id"`
	VpcName                  string        `mapstructure:"alicloud_vpc_name"`
	CidrBlock                string        `mapstructure:"alicloud_vpc_cidr_block"`
	VSwitchId                string        `mapstructure:"alicloud_vswitch_id"`
	VSwitchName              string        `mapstructure:"alicloud_vswitch_id"`
	InstanceName             string        `mapstructure:"alicloud_instance_name"`
	InternetChargeType       string        `mapstructure:"alicloud_internet_charge_type"`
	InternetMaxBandwidthOut  int           `mapstructure:"alicloud_internet_max_bandwith_out"`
	TemporaryKeyPairName     string        `mapstructure:"alicloud_temporary_key_pair_name"`
	WindowsPasswordTimeout   time.Duration `mapstructure:"alicloud_windows_password_timeout"`

	// Communicator settings
	Comm           communicator.Config `mapstructure:",squash"`
	SSHKeyPairName string              `mapstructure:"ssh_keypair_name"`
	SSHPrivateIp   bool                `mapstructure:"ssh_private_ip"`
	PublicKey      string              `mapstructure:"ssh_private_key_file"`
}

func (c *RunConfig) Prepare(ctx *interpolate.Context) []error {
	if c.SSHKeyPairName == "" && c.TemporaryKeyPairName == "" &&
		c.Comm.SSHPrivateKey == "" && c.Comm.SSHPassword == "" {

		c.TemporaryKeyPairName = fmt.Sprintf("packer_%s", uuid.TimeOrderedUUID())
	}

	if c.Comm.Type == "" {
		c.Comm.Type = default_comm_type
	}

	if c.Comm.SSHTimeout == 0 {
		c.Comm.SSHTimeout = ssh_time_out
	}

	if c.Comm.SSHPort == 0 {
		c.Comm.SSHPort = default_port
	}

	// Validation
	errs := c.Comm.Prepare(ctx)
	if c.AlicloudSourceImage == "" {
		errs = append(errs, errors.New("A alicloud_source_image must be specified"))
	}

	if c.InstanceType == "" {
		errs = append(errs, errors.New("An aliclod_instance_type must be specified"))
	}

	if c.UserData != "" && c.UserDataFile != "" {
		errs = append(errs, fmt.Errorf("Only one of alicloud_user_data or alicloud_user_data_file can be specified."))
	} else if c.UserDataFile != "" {
		if _, err := os.Stat(c.UserDataFile); err != nil {
			errs = append(errs, fmt.Errorf("alicloud_user_data_file not found: %s", c.UserDataFile))
		}
	}

	return errs
}
