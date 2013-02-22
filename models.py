import datetime
from flask import url_for
from extensions import db
from helpers import slugify
import json
from decimal import Decimal

from collections import OrderedDict


class Organization(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100))
    username = db.Column(db.String(50), unique=True)

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

    def __repr__(self):
        return '<User %s>' % self.login

