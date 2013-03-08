from fabric.api import *
from cuisine import *


def install_packages():
    package_ensure('build-essential')
    package_ensure('git')
    package_ensure('python-dev')
    package_ensure('python-setuptools')
    package_ensure('apache2')
    package_ensure('libapache2-mod-wsgi')
    package_ensure('postgresql-server-dev-9.1')
    package_ensure('postgresql-9.1')
    package_ensure('postgresql-client-9.1')
    package_ensure('vim')
    package_ensure('memcached')


def install_pip():
    sudo('easy_install -U pip')
    sudo('pip install --upgrade pip')


def sudo_if_no_dir(cmd, dir):
    if run('[ -d %s ] && echo 1 || echo 0' % dir) == '0':
        sudo(cmd)


def checkout():
    sudo_if_no_dir('git clone git://github.com/dlapiduz/govcode.git %s' % env.app_directory, env.app_directory)

def set_permissions():
    sudo('chown -R www-data:www-data /www')
    sudo('chmod -R g+w /www')
    sudo('usermod -G www-data ubuntu')

def stage():
    # environment variables
    env.app_directory = '/www/govcode.org'

    # let's get this server configured
    install_packages()
    install_pip()
    checkout()
    set_permissions()

    sudo('pip install -r %s/requirements.txt' % env.app_directory)
