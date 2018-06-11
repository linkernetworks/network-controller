# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  config.vm.hostname = 'dev'

  config.vm.provision "shell", privileged: false, inline: <<-SHELL
    set -e -x -u
    sudo apt-get update
    sudo apt-get install -y vim git build-essential openvswitch-switch tcpdump unzip tig
    # Env for proto
    PROTOC_RELEASE="https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip"
    PROTOC_TARGET="${HOME}/protoc"
    # Install Docker
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    sudo apt-get update
    sudo apt-get install -y docker-ce
    # Install Golang
    wget --quiet https://storage.googleapis.com/golang/go1.10.2.linux-amd64.tar.gz
    sudo tar -zxf go1.10.2.linux-amd64.tar.gz -C /usr/local/
    echo 'export GOROOT=/usr/local/go' >>  /home/$USER/.bashrc
    echo 'export GOPATH=$HOME/go' >> /home/$USER/.bashrc
    echo 'export PATH=/home/$USER/protoc/bin:$PATH:$GOROOT/bin:$GOPATH/bin' >> /home/$USER/.bashrc
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export PATH=/home/$USER/protoc/bin:$PATH:$GOROOT/bin:$GOPATH/bin
    # setup golang dir
    mkdir -p /home/$USER/go/src
    rm -rf /home/$USER/go1.9.1.linux-amd64.tar.gz
    # Download ovs CNI source
    git clone https://github.com/linkernetworks/network-controller go/src/github.com/linkernetworks/network-controller
    go get -u github.com/kardianos/govendor
    cd ~/go/src/github.com/linkernetworks/network-controller
    govendor sync
    # install protoc
    if [ ! -d "${PROTOC_TARGET}" ]; then curl -fsSL "$PROTOC_RELEASE" > "${PROTOC_TARGET}.zip"; fi
    if [ -f "${PROTOC_TARGET}.zip" ]; then unzip "${PROTOC_TARGET}.zip" -d "${PROTOC_TARGET}"; fi
    go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
  SHELL

  config.vm.provider :virtualbox do |v|
      v.customize ["modifyvm", :id, "--cpus", 2]
      v.customize ["modifyvm", :id, "--memory", 1024]
      v.customize ['modifyvm', :id, '--nicpromisc1', 'allow-all']
  end
end
