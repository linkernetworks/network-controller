# -*- mode: ruby -*- # vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  config.vm.hostname = 'network-controller-dev'
  config.vm.define vm_name = 'network_controller'
  config.vm.provision "file", source: "vagrant-docker.conf", destination: "/tmp/override.conf"

  config.vm.provision "shell", privileged: false, inline: <<-SHELL
    set -e -x -u
    sudo mkdir -p "/etc/systemd/system/docker.service.d/"
    sudo cp "/tmp/override.conf" "/etc/systemd/system/docker.service.d/override.conf"
    sudo apt-get update
    sudo apt-get install -y vim git build-essential openvswitch-switch tcpdump unzip tig
    # Env for proto
    PROTOC_RELEASE="https://github.com/google/protobuf/releases/download/v3.6.0/protoc-3.6.0-linux-x86_64.zip"
    PROTOC_TARGET="${HOME}/protoc"

    # Install Docker
    # kubernetes official max validated version: 17.03.2~ce-0~ubuntu-xenial
    export DOCKER_VERSION="17.06.2~ce-0~ubuntu"
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    sudo apt-get update
    sudo apt-get install -y docker-ce=${DOCKER_VERSION}

    # Install Kubernetes
    export KUBE_VERSION="1.11.0"
    export NET_IF_NAME="enp0s8"
    sudo apt-get install -y apt-transport-https curl
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
    echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee --append /etc/apt/sources.list.d/kubernetes.list
    sudo apt-get update
    sudo apt-get install -y kubectl kubelet=${KUBE_VERSION}-00 kubeadm=${KUBE_VERSION}-00

    # Disable swap
    sudo swapoff -a && sudo sysctl -w vm.swappiness=0
    sudo sed '/swap.img/d' -i /etc/fstab

    sudo kubeadm init --kubernetes-version v${KUBE_VERSION} --apiserver-advertise-address=172.17.8.100 --pod-network-cidr=10.244.0.0/16
    mkdir -p $HOME/.kube
    sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/config

    # Should give flannel the real network interface name
    wget --quiet https://raw.githubusercontent.com/coreos/flannel/v0.9.1/Documentation/kube-flannel.yml -O /tmp/kube-flannel.yml
    sed -i -- 's/"--kube-subnet-mgr"/"--kube-subnet-mgr", "--iface='"$NET_IF_NAME"'"/g' /tmp/kube-flannel.yml
    kubectl apply -f /tmp/kube-flannel.yml

    kubectl taint nodes --all node-role.kubernetes.io/master-

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
    # Download network-controller source
    git clone https://github.com/linkernetworks/network-controller go/src/github.com/linkernetworks/network-controller
    go get -u github.com/kardianos/govendor
    cd ~/go/src/github.com/linkernetworks/network-controller
    govendor sync
    # install protoc
    if [ ! -d "${PROTOC_TARGET}" ]; then curl -fsSL "$PROTOC_RELEASE" > "${PROTOC_TARGET}.zip"; fi
    if [ -f "${PROTOC_TARGET}.zip" ]; then unzip "${PROTOC_TARGET}.zip" -d "${PROTOC_TARGET}"; fi
    go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
  SHELL

  config.vm.network :private_network, ip: "172.17.9.100"
  config.vm.provider :virtualbox do |v|
      v.customize ["modifyvm", :id, "--cpus", 2]
      v.customize ["modifyvm", :id, "--memory", 4096]
      v.customize ['modifyvm', :id, '--nicpromisc1', 'allow-all']
  end
end
