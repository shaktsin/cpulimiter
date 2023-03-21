FROM alpine
WORKDIR app 

RUN apk add --no-cache bash gawk sed grep bc coreutils


COPY cpulim cpulim
COPY testapp testapp

CMD ["/bin/bash"]