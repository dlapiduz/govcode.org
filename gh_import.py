import requests
from pygithub3 import Github
from models import Organization, Repository, Commit, User
from extensions import db
import os
from flask.ext.script import Command


class GhImport(Command):

    gh = Github(token=os.getenv('GH_TOKEN'))

    def run(self):
        self.get_orgs()
        for org in Organization.query.filter_by(ignore=False).all():
            self.get_repos(org)
            print org
        for repo in Repository.query.all():
            
            try:
                print repo
                self.get_commits(repo)
            except Exception, e:
                db.session.rollback()
                print 'Error with ' + str(repo) + ": " + str(e)
                

        print "Finished Importing"


    def get_orgs(self):
        url = 'http://registry.usa.gov/accounts.json?service_id=github'

        r = requests.get(url)
        if r.status_code == 200:
            for acct in r.json['accounts']:
                org = Organization.query.filter_by(username=acct['account']).first()
                if org is None:
                    org = Organization()
                    org.name = acct['organization']
                    org.username = acct['account']
                    db.session.add(org)
            db.session.commit()
        else:
            print 'Error importing organizations'

    def get_repos(self, org):
        try:
            repos = self.gh.repos.list_by_org(org.username).all()
            for repo in repos:
                if not repo.fork:
                    r = Repository.query.filter_by(gh_id=repo.id).first()
                    if r is None:
                        r = Repository()
                    r.organization = org
                    r.gh_id = repo.id
                    r.name = repo.name
                    r.description = repo.description
                    r.forks = repo.forks
                    r.watchers = repo.watchers
                    r.size = repo.size
                    r.open_issues = repo.open_issues
                    db.session.add(r)
            db.session.commit()
        except:
            print 'error ' + org.name

    def get_commits(self, repo):
        last_commit = repo.commits.order_by(Commit.date.desc()).first()
        commit_pages = self.gh.repos.commits.list(user=repo.organization.username,
                                            repo=repo.name,
                                            sha="master")
        for commit in commit_pages.iterator():
            if last_commit and commit.sha == last_commit.sha:
                print 'Next'
                break
            c = Commit.query.filter_by(sha=commit.sha).first()
            if c is None:
                c = Commit()
                c.sha = commit.sha
                c.message = commit.commit.message
                c.repository = repo
                c.date = commit.commit.author.date
                if commit.author:
                    c.user = self.get_or_create_user(commit.author)
                print 'New commit: ' + c.sha
                db.session.add(c)
                db.session.commit()



    def get_or_create_user(self, user):
        u = User.query.filter_by(gh_id=user.id).first()
        if u is None:
            u = User()
            u.gh_id = user.id
            u.login = user.login
            u.avatar_url = user.avatar_url
            db.session.add(u)
            db.session.commit()
        return u
