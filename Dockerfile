FROM scratch

COPY /dns-lookup-server /dns-lookup-server

ENTRYPOINT [ "/dns-lookup-server" ]
CMD [ "" ]
