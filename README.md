# Mockbuster Movie API

![logo](./mockbuster-logo.jpeg)


## endpoints
```
/films                    [GET]    // list all films possible queryParams ? title=%s & rating=%s & category=%s  
/films/{filmId}           [GET]    // get film details by film
/films/{filmId}/comments  [GET]    // get comments for a film
/films/{filmId}/comments  [POST]   // add comment to film
```

## .env file
```dotenv
DB_HOST=localhost
DB_PORT=5555
DB_USER=postgres
DB_PASS=postgres
DB_NAME=dvdrental
API_DEBUG=true
```

