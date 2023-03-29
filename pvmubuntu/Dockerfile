FROM ubuntu:20.04 AS ubuntu-base

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update

RUN apt install -y pvm pvm-dev openssl
RUN rm -rf /var/lib/apt/lists/*
RUN apt clean

RUN useradd -rm -d /home/pvm -s /bin/bash -g root -G sudo pvm -p "$(openssl passwd -1 pvm)"
USER pvm
WORKDIR /home/pvm

FROM ubuntu-base AS ubuntu-with-sshd
USER root

RUN apt-get -qq update \
    && apt-get -qq --no-install-recommends install vim-tiny=2:8.1.* \
    && apt-get -qq --no-install-recommends install sudo=1.8.* \
    && apt-get -qq --no-install-recommends install python3-pip=20.0.* \
    && apt-get -qq --no-install-recommends install openssh-server=1:8.* \
    && apt-get -qq clean    \
    && rm -rf /var/lib/apt/lists/*

RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
RUN mkdir /var/run/sshd
RUN bash -c 'install -m755 <(printf "#!/bin/sh\nexit 0") /usr/sbin/policy-rc.d'
RUN ex +'%s/^#\zeListenAddress/\1/g' -scwq /etc/ssh/sshd_config
RUN ex +'%s/^#\zeHostKey .*ssh_host_.*_key/\1/g' -scwq /etc/ssh/sshd_config
RUN RUNLEVEL=1 dpkg-reconfigure openssh-server
RUN ssh-keygen -A -v
RUN update-rc.d ssh defaults

RUN ex +"%s/^%sudo.*$/%sudo ALL=(ALL:ALL) NOPASSWD:ALL/g" -scwq! /etc/sudoers

USER pvm
RUN ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519
#COPY --chown=pvm:root "./files/authorized_keys" /home/pvm/.ssh/authorized_keys

# Setup default command and/or parameters.
EXPOSE 22
CMD ["/usr/bin/sudo", "/usr/sbin/sshd", "-D", "-o", "ListenAddress=0.0.0.0"]