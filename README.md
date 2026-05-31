# Search engine
A simple inverted-index search engine implemented in Go using Redis as the storage backend.

This project exposes a small HTTP API that:
- Registers a URL and indexes all words found in its HTML;
- Returns URLs relevant to a user search query.

The indexing process includes tokenization, lemmatization, and stop-word filtering.

## Dependencies
- Golem for lemmatization;
- Goquery for HTML scraping;
- Redis for cache and key/value data.

# Redis config

Set all these enviroment variables to set up redis.
```
REDIS_SERVER=
REDIS_PASSWORD=
REDIS_PROTOCOL=
REDIS_DB=
```

## API Endpoints
### POST add/
Requires a url and indexes its content.

```
{
    "url": "https://www.google.com"
}
```

Behaviour:
1. The server fetches the URL;
2. HTML is tokenized and filtered;
3. For each token, a redis entry is created:

``
WORD -> SET(url1, url2) 
``

**Responses**:
- `201 OK` - content indexed.
- `4xx, 500` response from the HTTP request internal processing.

### POST /query
Searches for URLs relevant to the query.

Request body:
```
{
    "query": "how to study more efficiently"
}
```

Behaviour:
1. The query is tokenized and lemmatized;
2. For each token, Redis is queried for matching URLs;
3. The result is the reunion of all URL sets.

**Response**:
```
{
  "results": [
      "https://example.article.one/",
      "https://example.article.two/"
    ]
}
```
