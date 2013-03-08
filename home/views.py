from flask import Blueprint, request, redirect, render_template, url_for
from extensions import cache
from sqlalchemy import desc
from models import Organization, Repository, User


home = Blueprint('home', __name__, template_folder='templates')


@home.route('/')
@cache.cached(timeout=120)
def index():
    repositories = Repository.query.order_by('-forks').all()
    organizations = Organization.query.filter_by(ignore=False).order_by('name').all()

    context = {
        'repositories': repositories,
        'organizations': organizations
    }

    return render_template('home/index.html', **context)

@home.route('/repo/<slug>')
@cache.cached(timeout=120)
def repo(slug):
    repository = Repository.query.filter_by(id=slug).first()

    return render_template('home/repo.html', repository=repository)

@home.route('/contributors')
@cache.cached(timeout=120)
def contributors():
    users = User.query.order_by(desc(User.commit_count)).all()
    return render_template('home/contributors.html', users=users)

@home.route('/about')
@cache.cached(timeout=120)
def about():
    return render_template('home/about.html')