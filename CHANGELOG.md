# Change Log

## [v0.4.0](https://github.com/linkernetworks/network-controller/tree/v0.4.0) (2018-07-24)
[Full Changelog](https://github.com/linkernetworks/network-controller/compare/v0.3.0...v0.4.0)

**Merged pull requests:**

- \[Release\] bump to v0.4.0 [\#66](https://github.com/linkernetworks/network-controller/pull/66) ([John-Lin](https://github.com/John-Lin))
- Add add route feature into client [\#65](https://github.com/linkernetworks/network-controller/pull/65) ([sufuf3](https://github.com/sufuf3))
- Fixe golint & gofmt [\#64](https://github.com/linkernetworks/network-controller/pull/64) ([sufuf3](https://github.com/sufuf3))
- add git tag target [\#63](https://github.com/linkernetworks/network-controller/pull/63) ([John-Lin](https://github.com/John-Lin))
- Update changelog [\#62](https://github.com/linkernetworks/network-controller/pull/62) ([chenyunchen](https://github.com/chenyunchen))
- Add add route func in container [\#59](https://github.com/linkernetworks/network-controller/pull/59) ([sufuf3](https://github.com/sufuf3))

## [v0.3.0](https://github.com/linkernetworks/network-controller/tree/v0.3.0) (2018-07-19)
[Full Changelog](https://github.com/linkernetworks/network-controller/compare/v0.2.1...v0.3.0)

**Merged pull requests:**

- refactoring response message [\#61](https://github.com/linkernetworks/network-controller/pull/61) ([John-Lin](https://github.com/John-Lin))
- use CIDR for variable name [\#60](https://github.com/linkernetworks/network-controller/pull/60) ([John-Lin](https://github.com/John-Lin))
- VX-190 support setport vlan tag [\#58](https://github.com/linkernetworks/network-controller/pull/58) ([John-Lin](https://github.com/John-Lin))
- add tests [\#57](https://github.com/linkernetworks/network-controller/pull/57) ([John-Lin](https://github.com/John-Lin))

## [v0.2.1](https://github.com/linkernetworks/network-controller/tree/v0.2.1) (2018-07-06)
[Full Changelog](https://github.com/linkernetworks/network-controller/compare/v0.2.0...v0.2.1)

**Merged pull requests:**

- downgrade to protoc3.5.1 [\#56](https://github.com/linkernetworks/network-controller/pull/56) ([John-Lin](https://github.com/John-Lin))
- remove all vagrant stuff [\#55](https://github.com/linkernetworks/network-controller/pull/55) ([John-Lin](https://github.com/John-Lin))
- Add change log [\#52](https://github.com/linkernetworks/network-controller/pull/52) ([chenyunchen](https://github.com/chenyunchen))

## [v0.2.0](https://github.com/linkernetworks/network-controller/tree/v0.2.0) (2018-07-05)
[Full Changelog](https://github.com/linkernetworks/network-controller/compare/v0.1.1...v0.2.0)

**Merged pull requests:**

- update vagrantfile [\#54](https://github.com/linkernetworks/network-controller/pull/54) ([John-Lin](https://github.com/John-Lin))
- update protobuf to 3.6 [\#53](https://github.com/linkernetworks/network-controller/pull/53) ([John-Lin](https://github.com/John-Lin))
- Add port tag or vlan,trunk options setting [\#51](https://github.com/linkernetworks/network-controller/pull/51) ([chenyunchen](https://github.com/chenyunchen))
- functions for add dpdk port [\#50](https://github.com/linkernetworks/network-controller/pull/50) ([John-Lin](https://github.com/John-Lin))
- create bridge with type [\#49](https://github.com/linkernetworks/network-controller/pull/49) ([John-Lin](https://github.com/John-Lin))

## [v0.1.1](https://github.com/linkernetworks/network-controller/tree/v0.1.1) (2018-06-25)
[Full Changelog](https://github.com/linkernetworks/network-controller/compare/v0.1...v0.1.1)

**Merged pull requests:**

- Add bridge message [\#48](https://github.com/linkernetworks/network-controller/pull/48) ([chenyunchen](https://github.com/chenyunchen))
- Add message.pb.go for vortex to import [\#47](https://github.com/linkernetworks/network-controller/pull/47) ([chenyunchen](https://github.com/chenyunchen))
- Add openvswitch db.sock into tcp-daemonset [\#46](https://github.com/linkernetworks/network-controller/pull/46) ([sufuf3](https://github.com/sufuf3))
- add slack notification [\#45](https://github.com/linkernetworks/network-controller/pull/45) ([hwchiu](https://github.com/hwchiu))
- Add using the vagrant description [\#44](https://github.com/linkernetworks/network-controller/pull/44) ([sufuf3](https://github.com/sufuf3))
- Add Kubernetes into Vagrantfile [\#43](https://github.com/linkernetworks/network-controller/pull/43) ([sufuf3](https://github.com/sufuf3))
- Fix gofmt & golint errors [\#42](https://github.com/linkernetworks/network-controller/pull/42) ([sufuf3](https://github.com/sufuf3))

## [v0.1](https://github.com/linkernetworks/network-controller/tree/v0.1) (2018-06-13)
**Implemented enhancements:**

- Fix golint warning  [\#34](https://github.com/linkernetworks/network-controller/issues/34)
- Remove the unix file before listening if the server is running as unix server [\#32](https://github.com/linkernetworks/network-controller/issues/32)
- add go vet to do static check [\#6](https://github.com/linkernetworks/network-controller/pull/6) ([John-Lin](https://github.com/John-Lin))
- docker client & pause container find [\#5](https://github.com/linkernetworks/network-controller/pull/5) ([John-Lin](https://github.com/John-Lin))

**Fixed bugs:**

- Make the vagrantfile support docker-ce [\#28](https://github.com/linkernetworks/network-controller/issues/28)

**Merged pull requests:**

- Update the example yaml [\#41](https://github.com/linkernetworks/network-controller/pull/41) ([hwchiu](https://github.com/hwchiu))
- Add the deployment yaml files for client/server [\#39](https://github.com/linkernetworks/network-controller/pull/39) ([hwchiu](https://github.com/hwchiu))
- Fix golint errors [\#38](https://github.com/linkernetworks/network-controller/pull/38) ([sufuf3](https://github.com/sufuf3))
- Add client set container interface ip gateway [\#37](https://github.com/linkernetworks/network-controller/pull/37) ([chenyunchen](https://github.com/chenyunchen))
- Update clientpodtest.yaml [\#36](https://github.com/linkernetworks/network-controller/pull/36) ([John-Lin](https://github.com/John-Lin))
- Unlink the unix file before start grpc newserver [\#35](https://github.com/linkernetworks/network-controller/pull/35) ([sufuf3](https://github.com/sufuf3))
- add docker-ce [\#33](https://github.com/linkernetworks/network-controller/pull/33) ([John-Lin](https://github.com/John-Lin))
- VX-97 Add check IP Address format [\#31](https://github.com/linkernetworks/network-controller/pull/31) ([chenyunchen](https://github.com/chenyunchen))
- Install the openvswitch in the dockerifle. [\#30](https://github.com/linkernetworks/network-controller/pull/30) ([hwchiu](https://github.com/hwchiu))
- Remove unused function [\#29](https://github.com/linkernetworks/network-controller/pull/29) ([hwchiu](https://github.com/hwchiu))
- Add the netlink\_event\_tracker to handle netlink event [\#27](https://github.com/linkernetworks/network-controller/pull/27) ([hwchiu](https://github.com/hwchiu))
- VX-84 Add port statistics [\#26](https://github.com/linkernetworks/network-controller/pull/26) ([chenyunchen](https://github.com/chenyunchen))
- Add server daemonset yaml template files [\#25](https://github.com/linkernetworks/network-controller/pull/25) ([sufuf3](https://github.com/sufuf3))
- Support the TCP/UNIX parameter for server/client [\#24](https://github.com/linkernetworks/network-controller/pull/24) ([hwchiu](https://github.com/hwchiu))
- fix cni version at 0.6.0 [\#23](https://github.com/linkernetworks/network-controller/pull/23) ([John-Lin](https://github.com/John-Lin))
- Use net.ParseIP [\#22](https://github.com/linkernetworks/network-controller/pull/22) ([John-Lin](https://github.com/John-Lin))
- Fix OVS Interface Name [\#21](https://github.com/linkernetworks/network-controller/pull/21) ([chenyunchen](https://github.com/chenyunchen))
- fix error handle when docker get error [\#20](https://github.com/linkernetworks/network-controller/pull/20) ([John-Lin](https://github.com/John-Lin))
- Add the testing for docker package [\#18](https://github.com/linkernetworks/network-controller/pull/18) ([hwchiu](https://github.com/hwchiu))
- Add configureiface - VX-84 [\#17](https://github.com/linkernetworks/network-controller/pull/17) ([chenyunchen](https://github.com/chenyunchen))
- Check container ID if empty [\#16](https://github.com/linkernetworks/network-controller/pull/16) ([chenyunchen](https://github.com/chenyunchen))
- Add Client \(VX-83\) [\#15](https://github.com/linkernetworks/network-controller/pull/15) ([sufuf3](https://github.com/sufuf3))
- Refactor the openvswitch directory [\#14](https://github.com/linkernetworks/network-controller/pull/14) ([hwchiu](https://github.com/hwchiu))
- Create server handler include findNetworkNamespacePath and connectBridge [\#13](https://github.com/linkernetworks/network-controller/pull/13) ([chenyunchen](https://github.com/chenyunchen))
- Fix the vet and add the vet in travis [\#12](https://github.com/linkernetworks/network-controller/pull/12) ([hwchiu](https://github.com/hwchiu))
- support add veth pairs [\#10](https://github.com/linkernetworks/network-controller/pull/10) ([John-Lin](https://github.com/John-Lin))
- VX-49 Add Interrupt Signal [\#9](https://github.com/linkernetworks/network-controller/pull/9) ([chenyunchen](https://github.com/chenyunchen))
- Add more testing [\#8](https://github.com/linkernetworks/network-controller/pull/8) ([hwchiu](https://github.com/hwchiu))
- rename package name to openvswitch [\#7](https://github.com/linkernetworks/network-controller/pull/7) ([John-Lin](https://github.com/John-Lin))
- \[WIP\] Alex/modify flow [\#4](https://github.com/linkernetworks/network-controller/pull/4) ([chenyunchen](https://github.com/chenyunchen))
- Add testing [\#3](https://github.com/linkernetworks/network-controller/pull/3) ([hwchiu](https://github.com/hwchiu))
- Add the vagrant file [\#2](https://github.com/linkernetworks/network-controller/pull/2) ([hwchiu](https://github.com/hwchiu))
- Hwchiu/add travis [\#1](https://github.com/linkernetworks/network-controller/pull/1) ([hwchiu](https://github.com/hwchiu))
