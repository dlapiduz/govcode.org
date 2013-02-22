from flask import Blueprint, request, redirect, render_template, url_for

from models import Organization, Repository



home = Blueprint('home', __name__, template_folder='templates')


@home.route('/', methods=['GET', 'POST'])
def index():
    repositories = Repository.query.order_by('-forks').all()

    return render_template('home/index.html', repositories=repositories)