# Welcome to TweetNotes

## Getting Started

Tweetnotes is a web application for analysing tweets manually. It fetches tweets of specified people and lets you take note for each tweet.

## How to run

### Via Docker Compose

Copy the environment variable file and fill in the credentials:

    cp .env.web.sample .env.web
    
    vim .env.web

Then run:

    docker-compose up
    

Then see http://192.168.59.103:9000/

### Locally

* Make sure you have golang 1.4.2+ installed. Then install dependencies. 


    go get github.com/revel/revel
    
    go get gopkg.in/mgo.v2
    
    go get github.com/ChimeraCoder/anaconda
    
    go get github.com/revel/cmd/revel
    
    
* Make sure MongoDB instance is running on localhost:27071. 

 * Then set the environment variables:


    export CONSUMER_KEY=abcd...
    
    export CONSUMER_SECRET=efgh...
    
    export ACCESS_TOKEN=1234...
    
    export ACCESS_TOKEN_SECRET=6789....
    
    export MONGODB_DBNAME=tweetnotes


* Then start the server:


    revel run github.com/aladagemre/tweetnotes
 

 
Then see http://localhost:9000/

## Screenshots

Homepage:

![Screenshot](https://github.com/aladagemre/tweetnotes/blob/master/screenshot0.png)

Fill in the username and press the button and you will see:

![Screenshot](https://github.com/aladagemre/tweetnotes/blob/master/screenshot1.png)


