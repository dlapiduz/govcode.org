Govcode - Government Open Source Projects
=============

## What is this?

Govcode is an application that lists government open source projects.
The purpose is to track what is being worked on and build analytics on top of it.

## How do I run it?


1. Clone this repo or a fork of it
   ```bash
   git clone https://github.com/cfpb/collab.git
   cd collab
   ```

1. Make sure you have `pip` and `virtualenv` installed:

   - http://pip.readthedocs.org/en/latest/installing.html

   - http://www.virtualenv.org/en/latest/virtualenv.html#installation

1. Edit the config settings locally to match your environment.
   You can find them at [/config/__init__.py](https://github.com/dlapiduz/govcode/blob/master/config/__init__.py)

1. Create the database tables:

   ```bash
   python ./manage.py shell
   >>> from extensions import db
   >>> db.create_all()
   ```

1. Run the Flask server:

   ```bash
   python ./manage.py runserver
   ```

1. Go to <http://localhost:5000> in your browser.

## How to import the data?

At this point your database should be empty.
To actually import the packages you need to run the import script.

1. First you need to get a Github:

   Go to https://github.com/settings/tokens/new, create a token with `public_repo` access and copy the token.

1. Store the token in an environment variable:

   ```bash
   export GH_TOKEN=xxxxxxxxxxxx
   ```

1. Run the import script:

   ```bash
   python gh_import.py
   ```
