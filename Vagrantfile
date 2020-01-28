Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/bionic64"
  config.disksize.size = '50GB'

  config.vm.network "forwarded_port", guest: 80, host: 80, host_ip: "127.0.0.1"
  config.vm.network "forwarded_port", guest: 9000, host: 9000, host_ip: "127.0.0.1"
  config.vm.network "forwarded_port", guest: 52000, host: 52000, host_ip: "127.0.0.1"

  config.vm.provider "virtualbox" do |vb|
    vb.cpus = "4"
    vb.memory = "12288"
 end

 config.vm.provision "shell", inline: <<-SHELL
     apt-get update
     apt-get upgrade -y
     apt-get install make unzip -y
     snap install go --classic

     curl -fsSL https://get.docker.com -o get-docker.sh
     sh get-docker.sh
     usermod -aG docker vagrant


     cd /vagrant
     make build -j
     make init
     make plan
     make apply
   SHELL
end
