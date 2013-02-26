#!/usr/bin/env python
# Set the path
import os, sys
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from flask.ext.script import Manager, Server
from app import create_app

app = create_app()

manager = Manager(app)

from gh_import import GhImport

# Turn on debugger by default and reloader
manager.add_command("runserver", Server(
    use_debugger=True,
    use_reloader=True,
    host='0.0.0.0')
)

manager.add_command("do_import", GhImport())


if __name__ == "__main__":
    manager.run()