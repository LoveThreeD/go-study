FROM alpine
ADD study /study
ENTRYPOINT [ "/study" ]
