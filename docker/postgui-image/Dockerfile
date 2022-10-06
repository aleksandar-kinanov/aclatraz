FROM node:16
ENV CI=true

RUN mkdir /opt/postgui && git clone https://github.com/priyank-purohit/PostGUI.git /opt/postgui

WORKDIR /opt/postgui

RUN rm -f package-lock.json && npm install

CMD [ "npm" , "run", "start"]