import datetime
from flask import url_for
from extensions import db
from helpers import slugify
import json
from decimal import Decimal
from sqlalchemy import event
from sqlalchemy import desc

from collections import OrderedDict


class Organization(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100))
    username = db.Column(db.String(50), unique=True)
    slug = db.Column(db.String(100), unique=True)
    ignore = db.Column(db.Boolean, default=False)

    def __repr__(self):
        return '<Organization %s>' % self.name

class Repository(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    gh_id = db.Column(db.Integer)
    name = db.Column(db.String(50))
    forks = db.Column(db.Integer)
    watchers = db.Column(db.Integer)
    size = db.Column(db.Integer)
    open_issues = db.Column(db.Integer)
    description = db.Column(db.Text)
    organization_id = db.Column(db.Integer, db.ForeignKey('organization.id'))
    organization = db.relationship('Organization',
        backref=db.backref('repositories', lazy='dynamic'))
    users = db.relationship("User", secondary="commit")
    slug = db.Column(db.String(100), unique=True)

    def last_commit(self):
        return self.commits.order_by(desc(Commit.date)).first()

    def __repr__(self):
        return '<Repository %s>' % self.name

class Commit(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    sha = db.Column(db.String(70))
    message = db.Column(db.Text)
    date = db.Column(db.DateTime)
    repository_id = db.Column(db.Integer, db.ForeignKey('repository.id'))
    repository = db.relationship('Repository',
        backref=db.backref('commits', lazy='dynamic'))
    user_id = db.Column(db.Integer, db.ForeignKey('user.id'))
    user = db.relationship('User',
        backref=db.backref('commits', lazy='dynamic'))

    def __repr__(self):
        return '<Commit %s>' % self.message

class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    gh_id = db.Column(db.Integer)
    login = db.Column(db.String(70))
    avatar_url = db.Column(db.String(255))
    slug = db.Column(db.String(70), unique=True)

    def __repr__(self):
        return '<User %s>' % self.login

def add_slug(mapper, connection, target):
    if isinstance(target, Organization):
        target.slug = slugify(target.name)
    elif isinstance(target, Repository):
        target.slug = slugify("-".join([target.organization.username,target.name]))
    elif isinstance(target, User):
        target.slug = slugify(target.login)

event.listen(Organization, 'before_insert', add_slug)
event.listen(Organization, 'before_update', add_slug)
event.listen(Repository, 'before_insert', add_slug)
event.listen(Repository, 'before_update', add_slug)
event.listen(User, 'before_insert', add_slug)
event.listen(User, 'before_update', add_slug)