import sys
sys.path.insert(0, '/var/www/govcode.org')

from app import create_app
from config import ProdConfig

application = app = create_app(ProdConfig)
