#Download base image ubuntu 16.04
FROM balenalib/rpi-raspbian 

RUN mkdir /root/.ssh/
COPY id_rsa /root/.ssh/
COPY id_rsa.pub /root/.ssh/

RUN chmod 600 /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa.pub

RUN  apt-get -yq update && apt-get -yqq install ssh

RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

COPY setup.sh /

RUN sh /setup.sh

CMD ["python", "/IndependentStudyAppilcationRepo/ProfCode/files/cameraServer.py"]

EXPOSE 8181
