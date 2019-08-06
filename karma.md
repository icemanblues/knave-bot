# Karma Bot

This is going to be a simple implementation of a Karma recording/ tracking  
bot for slack.

Everyone starts with 0 karma

## Increase/ Decrement Karma

* If someone says `@user ++` then increase their karma by 1
* If someone says `@user --` then decrease their karma by 1

Not sure if we want to accept additional `+` or `-`
for additional increment or decrement.  

`@user ++++` would increase karma by 3? and vice versa for `----`

Any additional text after the `++` is saved as a message.

No self-karma. Need to make sure that a user is not giving themselves karma.

What about `@channel` or `@here` or slack groups `@groups`
give that karma to everyone involved? Or distribute it around (rounding?)

## Get Karma scoreboard

* Retrieve my own karma
* List the top X karma on the leaderboard
* List the bottom X karma on the leaderboard
* Give karma to someone else
* help, list the karma commands

Commands

* `/karma ++ @user`
* `/karma -- @user`
* `/karma me`
* `/karma top`
* `/karma bot` or `/karma bottom`
* `/karma help`

## Karma Interactions Tractions

We should store each and every karma as a transaction (ledger). Then we can provide a monthly lookback as a report

`/karma report`

Perhaps even a year in review. Or something like a quarter in review

# API

## Hubot karma

* https://github.com/github/hubot-scripts/blob/master/src/scripts/karma.coffee

## Slack Events

* https://api.slack.com/events-api
* https://api.slack.com/events/api
  * `url_verification` https://api.slack.com/events/url_verification
  * `message.channels` https://api.slack.com/events/message.channels

## Data Persistance

Table for storing the individual interactions/ transactions

* from_user
* to_user
* channel
* slack workspace
* karma delta
* message
* created_at
* modified_at

Table for storing the summaries / total

* to_user
* slack workspace
* karma (total)
* created_at
* modified_at

Maybe another table for the individual interactions summarized into months or weeks

## RESTful API

```
Might not want `/workspace/:workspace/` in the path
Might be better suited as a required query param
```

GET `/karma/workspace/:workspace/user/:user`

returns karma total

POST `/karma/workspace/:workspace/user/:user`  

* ?command
* ?delta
* ?message

Updates/ Creates a transaction, to adjust karma

GET `/karma/workspace/:workspace/rankings/top`  
GET `/karma/workspace/:workspace/rankings/bot`  

* ?n=5

returns n users and their karma totals, in order.
(n defaults to 5)

GET `/karma/report`

* lots of query params
