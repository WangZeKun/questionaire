FROM ubuntu

WORKDIR /App
COPY ./questionaire /App/
COPY ./conf /App/conf

CMD [ "./questionaire" ]