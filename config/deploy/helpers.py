from fabric.api import sudo, cd
from cuisine import file_update

def run_as(user, cmd):
  return sudo("sudo su %s -c '%s'" % (user, cmd))


def postgresql_role_ensure(username,
                           password,
                           superuser=False,
                           createdb=False,
                           createrole=False,
                           inherit=True,
                           login=True):

    cmd = "psql -tAc \"SELECT 1 FROM pg_roles WHERE rolname = '\\''{}'\\''\""
    with cd('/'):
        result = run_as('postgres', cmd.format(username))
        if result != '1':
            opts = [
              'SUPERUSER' if superuser else 'NOSUPERUSER',
              'CREATEDB' if createdb else 'NOCREATEDB',
              'CREATEROLE' if createrole else 'NOCREATEROLE',
              'INHERIT' if inherit else 'NOINHERIT',
              'LOGIN' if login else 'NOLOGIN'
            ]
            sql = "CREATE ROLE {username} WITH {opts} PASSWORD '\\''{password}'\\''"
            sql = sql.format(username=username, opts=' '.join(opts), password=password)
            cmd = 'psql -U postgres -c "{0}"'.format(sql)
            run_as('postgres', cmd)

def file_update_line(file, startswith, update):
    def update_line(text):
        ret = []
        for line in text.split('\n'):
            if line.strip().startswith(startswith):
                ret.append(update)
            else:
                ret.append(line)
        return "\n".join(ret)
    return file_update(file, update_line)