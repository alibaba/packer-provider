## Alicloud [(Alibaba Cloud)](http://www.aliyun.com) packer provider

This is the official repository for the Alicloud packer provider.  
Currently it supports packer version ≥ v0.12.1.

If you are not planning to contribute to this repo, you can download the [compiled binaries](https://github.com/alibaba/packer-provider/releases) according to you platform, unzip and move 
them into the folder under the packer **PATH** such as **/usr/local/packer**.

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
    "access_key": "{{env `ALICLOUD_ACCESS_KEY`}}",
    "secret_key": "{{env `ALICLOUD_SECRET_KEY`}}"
  },
  "builders": [{
    "type":"alicloud",
    "access_key":"{{user `access_key`}}",
    "secret_key":"{{user `secret_key`}}",
    "region":"cn-beijing",
    "image_name":"packer_test2",
    "source_image":"centos_7_2_64_40G_base_20170222.vhd",
    "ssh_username":"root",
    "instance_type":"ecs.n1.tiny",
    "io_optimized":"true",
    "image_force_delete":"true"
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
    "access_key": "{{env `ALICLOUD_ACCESS_KEY`}}",
    "secret_key": "{{env `ALICLOUD_SECRET_KEY`}}"
  },
  "builders": [{
    "type":"alicloud",
    "access_key":"{{user `access_key`}}",
    "secret_key":"{{user `secret_key`}}",
    "region":"cn-beijing",
    "image_name":"packer_test",
    "source_image":"win2008_64_ent_r2_zh-cn_40G_alibase_20170118.vhd",
    "instance_type":"ecs.n1.tiny",
    "io_optimized":"true",
    "image_force_delete":"true",
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
### How to contribute code
* If you are not sure or have any doubts, feel free to ask and/or submit an issue or PR. We appreciate all contributions and don't want to create artificial obstacles that get in the way.
* Contributions are welcome and will be merged via PRs.

### Contributors
* zhuzhih2017(dongxiao.zzh@alibaba-inc.com)

### License
* This project is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/alibaba/packer-provider/blob/master/LICENSE) for the full license text.

### Refrence
* Pakcer document: https://www.packer.io/intro/

