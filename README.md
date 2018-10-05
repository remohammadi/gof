
# gof

Game of Fakes: An Entertaining Approach to Analyzing the Impact of Awareness on Fake News Content and Trust in News Media

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ go get -u github.com/remohammadi/gof
$ cd $GOPATH/src/github.com/remohammadi/gof
$ heroku local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```
