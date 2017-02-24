# packer-provider
##Install
- Download the correct packer from you platform from https://www.packer.io/downloads.html
- Install packer according to the guide from https://www.packer.io/docs/installation.html
- Install Go according to the guide from [https://golang.org/doc/install](https://golang.org/doc/install)
- Setup your access key and secret key in the environment variables according to platform, for example In Linux platform with default bash, open your .bashrc in your home directory and add following two lines<p>
    ```aidl
        export ALICLOUD_ACCESS_KEY="access key value"
        
        export ALICLOUD_SECRET_KEY="secret key value"
     ```
- Open a terminator and clone Alicloud packer provider and build,install and test<p>
  ```
  cd <$GOPATH>
  
  mkdir -p src/github.com/alibaba/
  
  cd <$GOPATH>/src/github.com/alibaba/
  
  git clone https://github.com/alibaba/packer-provider
  
  cd <$GOPTH>/src/github.com/alibaba/packer-provider
    
  make all
  
  sorce ~/.bashrc
  
  packer build example/alicloud.json
  ```
 If output similar as following, configurations, you can now start the journey of alicloud with packer support
 ```
    alicloud output will be in this color.
    
    ==> alicloud: Force delete flag found, skipping prevalidating Alicloud ECS Image Name
        alicloud: Found Image ID: centos7u2_64_40G_cloudinit_20160728.raw
    ==> alicloud: allocated eip address 121.196.193.14
    ==> alicloud: Instance starting
    ==> alicloud: Waiting for SSH to become available...
    ==> alicloud: This machine's host=121.196.193.14
    ==> alicloud: This machine's host=121.196.193.14
    ==> alicloud: This machine's host=121.196.193.14
    ==> alicloud: Connected to SSH!
    ==> alicloud: Provisioning with shell script: /var/folders/3q/w38xx_js6cl6k5mwkrqsnw7w0000gn/T/packer-shell170579778
```
##Example
###Create a simple image with redis installed
```
{
  "variables": {
    "alicloud_access_key": "{{env `ALICLOUD_ACCESS_KEY`}}",
    "alicloud_secret_key": "{{env `ALICLOUD_SECRET_KEY`}}"
  },
  "builders": [{
    "type":"alicloud",
    "alicloud_access_key":"{{user `alicloud_access_key`}}",
    "alicloud_secret_key":"{{user `alicloud_secret_key`}}",
    "alicloud_region":"cn-beijing",
    "alicloud_image_name":"packer_test2",
    "alicloud_source_image":"centos7u2_64_40G_cloudinit_20160728.raw",
    "ssh_username":"root",
    "alicloud_instance_type":"ecs.n1.tiny",
    "alicloud_io_optimized":"true",
    "alicloud_image_force_delete":"true"
  }],
  "provisioners": [{
    "type": "shell",
    "inline": [
      "sleep 30",
      "yum install redis.x86_64 -y"
    ]
  }]
}

```
###Create a simple image for windows
```aidl
{
  "variables": {
    "alicloud_access_key": "{{env `ALICLOUD_ACCESS_KEY`}}",
    "alicloud_secret_key": "{{env `ALICLOUD_SECRET_KEY`}}"
  },
  "builders": [{
    "type":"alicloud",
    "alicloud_access_key":"{{user `alicloud_access_key`}}",
    "alicloud_secret_key":"{{user `alicloud_secret_key`}}",
    "alicloud_region":"cn-beijing",
    "alicloud_image_name":"packer_test",
    "alicloud_source_image":"win2008_64_ent_r2_zh-cn_40G_alibase_20170118.vhd",
    "alicloud_instance_type":"ecs.n1.tiny",
    "alicloud_io_optimized":"true",
    "alicloud_image_force_delete":"true",
    "communicator": "winrm",
    "winrm_port": 5985,
    "winrm_username": "Administrator",
    "winrm_password": "Test1234"
  }],
  "provisioners": [{
      "type": "powershell",
      "inline": ["dir c:\\"]
  }]
}

```
##
