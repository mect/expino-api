ARG arch
FROM multiarch/alpine:${arch}-edge

RUN mkdir /opt/kiosk

WORKDIR /opt/kiosk

COPY ./expino-backend /opt/kiosk/expino-backend
COPY ./frontend/build /opt/kiosk/expino-backend/frontend/build

EXPOSE 8080

CMD /opt/kioskexpino-backend