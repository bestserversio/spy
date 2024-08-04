FROM debian:12

# Update and install Go, Git, and others.
RUN apt update && apt install -y golang-1.19
RUN apt install -y git build-essential curl

# Add app user.
RUN groupadd app
RUN useradd -m -d /app -s /bin/bash -g app app

# Change work directory and user.
WORKDIR /app
USER app

# Clone, build, and install Spy.
RUN git clone --recursive https://github.com/bestserversio/spy.git

WORKDIR /app/spy
ENV PATH="/usr/lib/go-1.19/bin:${PATH}"

RUN ./build.sh

USER root
RUN make install

# Change back to app user.
WORKDIR /app
USER app
