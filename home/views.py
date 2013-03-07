from flask import Blueprint, request, redirect, render_template, url_for
from sqlalchemy import desc
from models import Organization, Repository, User


home = Blueprint('home', __name__, template_folder='templates')


@home.route('/')
def index():
    repositories = Repository.query.order_by('-forks').all()
    organizations = Organization.query.filter_by(ignore=False).order_by('name').all()

    context = {
        'repositories': repositories,
        'organizations': organizations
    }

    return render_template('home/index.html', **context)

@home.route('/repo/<slug>', methods=['GET', 'POST'])
def repo(slug):
    repository = Repository.query.filter_by(id=slug).first()

    return render_template('home/repo.html', repository=repository)