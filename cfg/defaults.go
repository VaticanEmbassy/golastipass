package cfg

var DefaultMapping string = `{
    "settings": {
        "index": {
            "refresh_interval": -1,
            "number_of_replicas": 0
        },
        "analysis": {
            "filter": {
                "tld_filter": {
                    "type": "pattern_capture",
                    "preserve_original": false,
                    "patterns": ["\\.([^\\.]+?)$"]
                }
            },
            "analyzer": {
                "lc_analyzer": {
                    "type": "custom",
                    "tokenizer": "keyword",
                    "filter": ["lowercase"]
                },
                "user_analyzer": {
                    "type": "custom",
                    "tokenizer": "user_tokenizer",
                    "filter": ["lowercase"]
                },
                "domain_analyzer": {
                    "type": "custom",
                    "tokenizer": "domain_tokenizer",
                    "filter": ["lowercase"]
                },
                "domain_notld_analyzer": {
                    "type": "custom",
                    "tokenizer": "domain_notld_tokenizer",
                    "filter": ["lowercase"]
                },
                "tld_analyzer": {
                    "type": "custom",
                    "tokenizer": "tld_tokenizer",
                    "filter": ["lowercase"]
                }
            },
            "tokenizer": {
                "user_tokenizer": {
                    "type": "pattern",
                    "pattern": "(.+?)@",
                    "group": 1
                },
                "domain_tokenizer": {
                    "type": "pattern",
                    "pattern": "@(.+)",
                    "group": 1
                },
                "domain_notld_tokenizer": {
                    "type": "pattern",
                    "pattern": "@(.+)\\.",
                    "group": 1
                },
                "tld_tokenizer": {
                    "type": "pattern",
                    "pattern": "\\.([^\\.]+?)$",
                    "group": 1
                }
            },
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
        "account": {
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
    }
}`
