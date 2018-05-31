# -*- mode: ruby -*-
# vi: set ft=ruby :

$script = <<SCRIPT
set -e -x -u
sudo apt-get -qq update
sudo apt-get -y -qq install -y vim git htop dh-autoreconf libssl-dev libcap-ng-dev build-essential openssl python python-pip openvswitch-switch
sudo pip install six

#### Install Golang
wget --quiet https://storage.googleapis.com/golang/go1.10.2.linux-amd64.tar.gz
sudo tar -zxf go1.10.2.linux-amd64.tar.gz -C /usr/local/
echo 'export GOROOT=/usr/local/go' >> /home/vagrant/.bashrc
echo 'export GOPATH=$HOME/go' >> /home/vagrant/.bashrc
echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> /home/vagrant/.bashrc
source /home/vagrant/.bashrc
mkdir -p /home/vagrant/go/src
rm -rf /home/vagrant/go1.10.2.linux-amd64.tar.gz

#### govendor
go get -u github.com/kardianos/govendor
SCRIPT

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  config.vm.hostname = "devstack"
  config.vm.provision "shell", privileged: false, inline: $script
  # config.vm.network "private_network", ip: "192.168.111.222"
  config.vm.network "public_network"

  config.vm.provider :virtualbox do |v|
      v.customize ["modifyvm", :id, "--cpus", 2]
      v.customize ["modifyvm", :id, "--memory", 1024]
      v.customize ['modifyvm', :id, '--nicpromisc2', 'allow-all']
  end # end provider
end
