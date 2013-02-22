# -*- coding: utf-8 -*-
"""
    Unit Tests
    ~~~~~~~~~~

    Define TestCase as base class for unit tests.
    Ref: http://packages.python.org/Flask-Testing/
"""

from flask.ext.testing import TestCase as Base, Twill

from app import create_app
from config import TestConfig
from extensions import db


class TestCase(Base):


    def create_app(self):
        """Create and return a testing flask app."""

        app = create_app(TestConfig)
        self.twill = Twill(app, port=3000)
        return app


    def _test_get_request(self, endpoint, template=None):
        response = self.client.get(endpoint)
        self.assert_200(response)
        # NOT Supported yet
        # if template:
        #     self.assertTemplateUsed(name=template)
        return response
