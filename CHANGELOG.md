## 1.3.1 (December 14, 2018)

MPROVEMENTS:

- Support `-force` option and copied images and snapshots if corresponding options are specified ([#57](https://github.com/alibaba/packer-provider/pull/57))

## 1.3.0 (November 26, 2018)

IMPROVEMENTS:

- Support wait_snapshot_ready_timeout for much bigger disk ([#53](https://github.com/alibaba/packer-provider/pull/53))
- Apply tags to relevant snapshots ([#54](https://github.com/alibaba/packer-provider/pull/54))
- Update windows examples ([#55](https://github.com/alibaba/packer-provider/pull/55))

## 1.2.5 (November 18, 2018)

IMPROVEMENTS:

- Support creating image without data disks ([#50](https://github.com/alibaba/packer-provider/pull/50))

## 1.2.4 (October 31, 2018)

IMPROVEMENTS:

- add options for system disk properties ([#48](https://github.com/alibaba/packer-provider/pull/48))

## 1.2.3 (September 28, 2018)

IMPROVEMENTS:

- Support disable_stop_instance option for some specific scenarios ([#45](https://github.com/alibaba/packer-provider/pull/45))

## 1.2.2 (September 16, 2018)

IMPROVEMENTS:

- Support adding tags to image ([#43](https://github.com/alibaba/packer-provider/pull/43))

## 1.2.1 (September 11, 2018)

IMPROVEMENTS:

- Support ssh with private ip address ([#42](https://github.com/alibaba/packer-provider/pull/42))

## 1.2.0 (August 14, 2018)

IMPROVEMENTS:

- Support describing marketplace image ([#39](https://github.com/alibaba/packer-provider/pull/39))
- Sync with official packer ([#38](https://github.com/alibaba/packer-provider/pull/38))

## 1.1.3 (August 3, 2017)

IMPROVEMENTS:

- add international site support ([#31](https://github.com/alibaba/packer-provider/pull/31))

## 1.1.2 (July 20, 2017)

IMPROVEMENTS:

- Refactor the code and enhance the retry logic to reduce the timeout failure.

## 1.1.1 (May 24, 2017)

BUG FIXES:

- Fix the missing parameter paytype when allocate eip

## 1.1 (March 12, 2017)

IMPROVEMENTS:

- Add local image import function ([#16](https://github.com/alibaba/packer-provider/pull/16))

## 1.0 (March 3, 2017)

IMPROVEMENTS:

- Add alicloud official announcement and orgnize the structure of samples ([#10](https://github.com/alibaba/packer-provider/pull/10))
