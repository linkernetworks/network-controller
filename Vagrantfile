# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  config.vm.hostname = 'network_controller-dev'
  config.vm.define vm_name = 'network_controller'
  config.vm.provision "file", source: "vagrant-docker.conf", destination: "/tmp/override.conf"

  config.vm.provision "shell", privileged: false, inline: <<-SHELL
    set -e -x -u
    sudo mkdir -p "/etc/systemd/system/docker.service.d/"
    sudo cp "/tmp/override.conf" "/etc/systemd/system/docker.service.d/override.conf"
    sudo apt-get update
    sudo apt-get install -y vim git build-essential openvswitch-switch tcpdump unzip tig
    # Env for proto
    PROTOC_RELEASE="https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip"
    PROTOC_TARGET="${HOME}/protoc"
    # Install Docker
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    sudo apt-get update
    sudo apt-get install -y docker-ce=17.03.2~ce-0~ubuntu-xenial
    # Install Kubernetes
    sudo apt-get install -y apt-transport-https curl
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
    echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee --append /etc/apt/sources.list.d/kubernetes.list
    sudo apt-get update
    sudo apt-get install -y kubelet kubeadm kubectl
    sudo swapoff -a
    sudo kubeadm init --apiserver-advertise-address=172.17.9.100 --pod-network-cidr=10.244.0.0/16
    mkdir -p $HOME/.kube
    sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/config
    kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/v0.9.1/Documentation/kube-flannel.yml
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
      v.customize ["modifyvm", :id, "--memory", 1024]
      v.customize ['modifyvm', :id, '--nicpromisc1', 'allow-all']
  end
end
