from fabric.api import *
from fabistrano import deploy

from config.deploy import setup

env.hosts = ["govcode.org"]
env.base_dir = '/www'
env.app_name = 'govcode'
env.git_clone = 'git://github.com/dlapiduz/govcode.git'
env.restart_cmd = 'kill -HUP `supervisorctl pid gunicorn`'
