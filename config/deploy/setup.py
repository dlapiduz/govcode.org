import os
from fabric.api import *
from cuisine import *
from helpers import *

def install_packages():
    packages = ['build-essential', 'git', 'python-dev', 
                'python-setuptools', 'nginx', 'nginx-common', 
                'postgresql-server-dev-9.1', 'postgresql-9.1', 
                'postgresql-client-9.1', 'vim', 'memcached', 
                'htop'  , 'supervisor', 'python-pip', 'libmemcached-dev',
                'zlib1g-dev', 'libssl-dev']


    for package in packages:
        package_ensure(package)

def setup_user():
    mode_sudo()
    user_ensure('diego', shell="/bin/bash")
    key = file_local_read('config/deploy/files/diego.pub')
    ssh_authorize('diego', key)
    file_update(
        "/etc/sudoers",
        lambda _:text_ensure_line(_, "diego ALL=NOPASSWD: ALL")
    )
    group_user_ensure('www-data', 'diego')

def setup_sshd():
    mode_sudo()
    file_update_line('/etc/ssh/sshd_config', 'PermitRootLogin', 'PermitRootLogin no')

def configure_database():
    mode_sudo()
    postgresql_role_ensure('govcode', 'govcode', createdb=True, login=True)

def setup_gunicorn():
    mode_sudo()
    python_package_ensure_pip('gunicorn')

def setup_supervisor():
    mode_sudo()
    file_ensure('/etc/supervisor/conf.d/govcode.conf', owner='root', group='root')
    text = file_local_read('config/deploy/files/supervisord.conf')
    text = text_template(text, dict(
                GH_TOKEN=os.environ['GH_TOKEN'],
                PG_CONN=os.environ['PG_REMOTE']))
    file_write('/etc/supervisor/conf.d/govcode.conf', text, owner='root', group='root')
    sudo('supervisorctl reread')
    sudo('supervisorctl update')
    sudo('supervisorctl start gunicorn')

@task
def setup_nginx():
    mode_sudo()
    file_upload('/etc/nginx/sites-available/govcode.org', 'config/deploy/files/nginx.conf')
    sudo('ln -s /etc/nginx/sites-available/govcode.org /etc/nginx/sites-enabled/govcode.org')
    sudo('service nginx restart')


@task(default=True)
def setup():
    # install_packages()
    # setup_user()
    # setup_sshd()
    # configure_database()
    # setup_gunicorn()
    # setup_nginx()
    setup_supervisor()

