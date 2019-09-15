package cfg

var DefaultMapping string = `{
    "settings": {
        "index": {
            "refresh_interval": -1,
            "number_of_replicas": 0
        },
        "analysis": {
            "normalizer": {
                "lc_normalizer": {
                    "type": "custom",
                    "char_filter": [],
                    "filter": ["lowercase"]
                }
            }
        }
    },
    "mappings": {
        "properties": {
            "email": {
                "type": "text",
                "analyzer": "simple",
                "fields": {
                    "raw": {
                        "type": "keyword",
                        "normalizer": "lc_normalizer"
                    }
                }
            },
            "username": {
                "type": "text",
                "analyzer": "simple",
                "fields": {
                    "raw": {
                        "type": "keyword",
                        "normalizer": "lc_normalizer"
                    }
                }
            },
            "domain": {
                "type": "keyword",
                "normalizer": "lc_normalizer"
            },
            "domain_notld": {
                "type": "keyword",
                "normalizer": "lc_normalizer"
            },
            "tld": {
                "type": "keyword",
                "normalizer": "lc_normalizer"
            },
            "password": {
                "type": "text",
                "analyzer": "simple",
                "fields": {
                    "raw": {
                        "type": "keyword"
                    }
                }
            },
            "password_length": {
                "type": "short"
            },
            "source": {
                "type": "short"
            }
        }
    }
}`
