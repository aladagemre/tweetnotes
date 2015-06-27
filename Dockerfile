FROM google/golang:1.4

#WORKDIR /gopath/src/app
#ADD . /gopath/src/app/
#RUN go get app


#ENTRYPOINT ["/gopath/bin/app"]

RUN go get github.com/revel/revel
RUN go get gopkg.in/mgo.v2
RUN go get github.com/ChimeraCoder/anaconda
RUN go get github.com/revel/cmd/revel

WORKDIR /gopath/src/github.com/aladagemre/tweetnotes
ADD . /gopath/src/github.com/aladagemre/tweetnotes
CMD ["revel", "run", "github.com/aladagemre/tweetnotes"]
EXPOSE 9000
#ENTRYPOINT ["/gopath"]
#CMD ["revel"]