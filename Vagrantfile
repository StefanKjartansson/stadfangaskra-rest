# -*- mode: ruby -*-
# vi: set ft=ruby :

BOX_NAME = "ubuntu"
BOX_URI = "http://files.vagrantup.com/precise64.box"

SCM, ORG, PROJECT = File.absolute_path(File.dirname(__FILE__)).split('/').last(3)
GOPATH = "/home/vagrant/go/src/#{SCM}/#{ORG}"
PPATH = "#{GOPATH}/#{PROJECT}"

Vagrant::Config.run do |config|
    config.vm.box = BOX_NAME
    config.vm.box_url = BOX_URI
    config.vm.network :hostonly, "33.33.33.18"
    config.vm.share_folder(PROJECT, PPATH, ".", :nfs => true)
    config.vm.provision :shell do |shell|
        shell.inline = <<-eos
            apt-get update -qq; 
            apt-get install -q -y python-software-properties nfs-common make git mercurial bzr;
            add-apt-repository -y ppa:gophers/go 
            apt-get update -qq;
            apt-get install -q -y golang-stable;
            su vagrant -c 'echo "export GOPATH=$HOME/go" >> $HOME/.bashrc'
            chown vagrant /home/vagrant/go
            chown vagrant /home/vagrant/go/src
            chown vagrant /home/vagrant/go/src/#{SCM}
            chown vagrant /home/vagrant/go/src/#{SCM}/#{ORG}
        eos
    end
end


Vagrant::VERSION >= "1.1.0" and Vagrant.configure("2") do |config|
    config.vm.provider :virtualbox do |vb|
        config.vm.box = BOX_NAME
        config.vm.box_url = BOX_URI
    end
end
