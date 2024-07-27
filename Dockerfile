FROM amd64/ubuntu

# The /app directory should act as the main application directory
WORKDIR /app

COPY . .

ARG CHROME_VERSION=125.0.6422.78

# install utils
RUN apt-get update \
  && apt-get install -y wget curl zip tar

# install chrome
# todo find version specific repo
# http://dl.google.com/linux/chrome/deb/pool/main/g/google-chrome-stable/google-chrome-stable_125.0.6422.112-1_amd64.deb
RUN curl -LO https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
RUN apt-get install -y ./google-chrome-stable_current_amd64.deb
RUN rm google-chrome-stable_current_amd64.deb

# install chromedriver
# import get_driver script
# get_driver CHROME_VERSION
# requires jq package
RUN wget https://storage.googleapis.com/chrome-for-testing-public/${CHROME_VERSION}/linux64/chromedriver-linux64.zip \
  && unzip chromedriver-linux64.zip \
  && chmod +x chromedriver-linux64/chromedriver \
  && mv chromedriver-linux64/chromedriver /usr/local/bin/chromedriver \
  && rm -dfr chromedriver_linux64.zip

# install go
# todo: change to go based image
RUN wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go \
  && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# Start the app using serve command
CMD [ "go", "test", "-count=1", "-v", "test/driver_test.go", "-run", "TestDriver" ]

