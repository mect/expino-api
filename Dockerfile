ARG arch
FROM multiarch/alpine:${arch}-edge

COPY ./expino-backend /expino-backend

CMD /expino-backend