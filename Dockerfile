FROM centos:8

WORKDIR /

RUN curl -s https://packagecloud.io/install/repositories/choria/release/script.rpm.sh | bash && \
    yum -y update && \
    yum -y install choria && \
    yum -y clean all

RUN groupadd --gid 2048 choria && \
    useradd -c "Choria Orchestrator - choria.io" -m --uid 2048 --gid 2048 choria && \
    chown -R choria:choria /etc/choria

USER choria

ENTRYPOINT ["/bin/choria"]