# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "https://cloud-images.ubuntu.com/vagrant/trusty/current/trusty-server-cloudimg-amd64-vagrant-disk1.box"

  config.vm.synced_folder ".", "/zk-backyp"
  #config.vm.synced_folder "../../../bosh-environment", "/bosh-env"

  $script = <<SCRIPT
#apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
#echo deb https://get.docker.io/ubuntu docker main > /etc/apt/sources.list.d/docker.list
apt-get update
apt-get install -y git golang monit zookeeper openjdk-7-jre

export GOPATH=~/.go/
echo export GOPATH=~/.go/ >> ~/.bashrc

cp /zk-backyp/test/monit            /etc/monit/conf.d/zookeeper
echo "set httpd port 2812 and"   >> /etc/monit/conf.d/httpd
echo "    use address localhost" >> /etc/monit/conf.d/httpd
echo "    allow localhost"       >> /etc/monit/conf.d/httpd
sed -i "s/set daemon 120/set daemon 2/g" /etc/monit/monitrc
service monit restart
SCRIPT

  config.vm.provision "shell", inline: $script

end
