import re
from unicodedata import normalize

def slugify(text, delim=u'-'):
    """Generates an ASCII-only slug."""
    punct_re = re.compile(r'[\t !"#$%&\'()*\-/<=>?@\[\\\]^_`{|},.]+')

    result = []
    for word in punct_re.split(text.lower()):
        if word:
            word = normalize('NFKD', unicode(word)).encode('ascii', 'ignore')
            result.append(word)
    return unicode(delim.join(result))

def es_search(term):
    conn = ES('127.0.0.1:9200')
    q = StringQuery(term)
    s = Search(q, start=0, size=50)
    return conn.search(s)
