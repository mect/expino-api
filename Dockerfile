ARG arch
FROM multiarch/alpine:${arch}-edge

RUN mkdir -p /opt/kiosk/expino-backend/frontend/build

WORKDIR /opt/kiosk

COPY ./expino-backend /opt/kiosk/expino-backend
COPY ./frontend/build /opt/kiosk/frontend/build

EXPOSE 8080

CMD /opt/kiosk/kioskexpino-backend