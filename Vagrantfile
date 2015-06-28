# -*- mode: ruby -*-
# vi: set ft=ruby :

ssh_user = ENV.fetch('USER', nil)

$script = <<-SCRIPT
echo "- updating deb repository"

cd /vagrant

echo "- installing golang 1.4.2"
if [ ! -d /opt/go ]; then
  apt-get update > /dev/null
  apt-get install -y --force-yes -qq git > /dev/null
  cd /tmp
  wget -q https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
  tar -xf go1.4.2.linux-amd64.tar.gz
  mv go /opt
  mkdir -p /opt/gopkg
  chown -R vagrant:vagrant /opt/go /opt/gopkg
fi

echo "- installing ELK"
if [ ! -d /usr/share/elasticsearch ]; then
  apt-get update > /dev/null
  apt-get install -y --force-yes -qq software-properties-common > /dev/null
  wget -qO - https://packages.elasticsearch.org/GPG-KEY-elasticsearch | sudo apt-key add - > /dev/null
  sudo add-apt-repository "deb http://packages.elasticsearch.org/elasticsearch/1.4/debian stable main"
  sudo add-apt-repository "deb http://packages.elasticsearch.org/logstash/1.4/debian stable main"

  apt-get update > /dev/null
  apt-get install -y --force-yes -qq elasticsearch logstash openjdk-7-jre
  update-rc.d elasticsearch defaults 95 10
  /etc/init.d/elasticsearch start
fi

if [ ! -f /usr/bin/redis-server ]; then
  apt-get update > /dev/null
  apt-get install -y --force-yes -qq redis-server > /dev/null
fi

if [ ! -f /usr/bin/carbon-cache ]; then
  echo "- installing graphite"
  apt-get update > /dev/null
  apt-get install -y --force-yes -qq graphite-carbon graphite-web libpq-dev python-psycopg2 python-memcache > /dev/null
  sed -i "/#SECRET_KEY = 'UNSAFE_DEFAULT'/c\SECRET_KEY = 'a_salty_string'" /etc/graphite/local_settings.py
  sed -i "s/America\\/Los_Angeles/UTC/" /etc/graphite/local_settings.py
  sed -i "s/#TIME_ZONE/TIME_ZONE/g" /etc/graphite/local_settings.py
  sed -i "/#USE_REMOTE_USER_AUTHENTICATION = True/c\USE_REMOTE_USER_AUTHENTICATION = True" /etc/graphite/local_settings.py
  sed -i "s/\\/var\\/lib\\/graphite\\/graphite.db/graphite/g" /etc/graphite/local_settings.py
  sed -i "s/django.db.backends.sqlite3/django.db.backends.postgresql_psycopg2/g" /etc/graphite/local_settings.py
  sed -i "s/USER': ''/USER': 'graphite'/g" /etc/graphite/local_settings.py
  sed -i "s/PASSWORD': ''/PASSWORD': 'password'/g" /etc/graphite/local_settings.py
  sed -i "s/HOST': ''/HOST': '127.0.0.1'/g" /etc/graphite/local_settings.py
  sed -i "s/CARBON_CACHE_ENABLED=false/CARBON_CACHE_ENABLED=true/g" /etc/default/graphite-carbon
  sed -i "s/ENABLE_LOGROTATION = False/ENABLE_LOGROTATION = True/g" /etc/carbon/carbon.conf
  sudo service carbon-cache start > /dev/null

  apt-get install -y --force-yes -qq postgresql libpq-dev python-psycopg2 > /dev/null
  su - postgres -c "psql -c \\"CREATE USER graphite WITH PASSWORD 'password';\\"" > /dev/null
  su - postgres -c "psql -c \\"CREATE DATABASE graphite WITH OWNER graphite;\\"" > /dev/null
  graphite-manage syncdb --noinput > /dev/null
  echo "from django.contrib.auth.models import User; User.objects.create_superuser('admin', 'mail@example.com', 'password')" | graphite-manage shell > /dev/null

  echo "- installing apache2"
  apt-get install -y --force-yes -qq apache2 libapache2-mod-wsgi > /dev/null
  a2dissite 000-default > /dev/null
  cp /usr/share/graphite-web/apache2-graphite.conf /etc/apache2/sites-available
  a2ensite apache2-graphite > /dev/null
  service apache2 reload > /dev/null
fi

if [ ! -L /opt/influxdb/influxd ]; then
  echo "- installing influxdb"
  cd /tmp
  wget -q http://influxdb.s3.amazonaws.com/influxdb_0.9.0_amd64.deb
  sudo dpkg -i influxdb_0.9.0_amd64.deb > /dev/null
  /etc/init.d/influxdb start > /dev/null
fi

echo "- ensuring environment file is up to date"
ENV_FILE="/etc/environment"
ENV_TEMP=`cat "${ENV_FILE}"`
ENV_TEMP=$(echo -e "${ENV_TEMP}" | sed "/^GOPATH=/ d")
ENV_TEMP=$(echo -e "${ENV_TEMP}" | sed "/^GOROOT=/ d")
ENV_TEMP=$(echo -e "${ENV_TEMP}" | sed "/^PATH=/ d")
ENV_TEMP=$(echo -e "${ENV_TEMP}" | sed "/^SSH_USER=/ d")
ENV_TEMP="${ENV_TEMP}\nGOPATH='/opt/gopkg'"
ENV_TEMP="${ENV_TEMP}\nGOROOT='/opt/go'"
ENV_TEMP="${ENV_TEMP}\nPATH='/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/opt/gopkg/bin:/opt/go/bin:/opt/influxdb:/opt/logstash/bin:/var/lib/kibana/bin'"
ENV_TEMP="${ENV_TEMP}\nSSH_USER='#{ssh_user}'"
echo "$ENV_TEMP" | sed '/^$/d' | sort > $ENV_FILE

if ! grep -q cd-to-directory "/home/vagrant/.bashrc"; then
  echo "- setting up auto chdir on ssh"
  echo "\n[ -n \\"\\$SSH_CONNECTION\\" ] && cd /opt/gopkg/src/github.com/josegonzalez/metricsd # cd-to-directory" >> "/home/vagrant/.bashrc"
fi

echo -e "\n- ALL CLEAR! SSH access via 'vagrant ssh'"
echo "- Virtual Machine IP:"
ifconfig | grep "inet " | grep -v 127 | grep -v "addr:10.0" | cut -d':' -f2 | cut -d' ' -f1
SCRIPT


VAGRANTFILE_API_VERSION = "2"
Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "chef/ubuntu-14.04"
  config.vm.provision :shell, inline: $script
  config.vm.network "private_network", type: "dhcp"
  config.vm.synced_folder ".", "/opt/gopkg/src/github.com/josegonzalez/metricsd"
end
