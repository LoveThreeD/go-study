FROM alpine
ADD sTest /sTest
ENTRYPOINT [ "/sTest" ]
