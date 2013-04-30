# -*- coding: utf-8 -*-
import os


class BaseConfig(object):

    # Get app root path
    # ../../configs/config.py
    _basedir = os.path.abspath(os.path.dirname(os.path.dirname(__file__)))

    PROJECT = "govcode"
    DEBUG = False
    TESTING = False

    SECRET_KEY = '\x8137i\xd6\xc4X\xf4Y\xcf\xb3\x95M\xe4n\xbf\xd5\xb3\x18}l\x81\x03\xb5'

class DevConfig(BaseConfig):
    DEBUG = True
    SQLALCHEMY_DATABASE_URI = 'postgresql://diego@localhost:5432/govcode_dev'
    CACHE_CONFIG = { 'CACHE_TYPE': 'null' }



class TestConfig(BaseConfig):
    TESTING = True
    CSRF_ENABLED = False

    SQLALCHEMY_DATABASE_URI = 'mysql://root@localhost/lmm_test'
    CACHE_CONFIG = { 'CACHE_TYPE': 'null' }


class ProdConfig(BaseConfig):
    SQLALCHEMY_DATABASE_URI = os.getenv('PG_CONN_STR')
    CACHE_CONFIG = { 'CACHE_TYPE': 'memcached' }
