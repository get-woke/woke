# [Choice] Go version: 1, 1.16, 1.15
ARG VARIANT=1.17
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Option] Install Node.js
ARG DEBIAN_FRONTEND=noninteractive
ARG INSTALL_NODE="true"
ARG NODE_VERSION="lts/*"
RUN if [ "${INSTALL_NODE}" = "true" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi \
    # Install Python
    && apt-get update \
    && apt-get install -y python3.6 python3-pip \
    # Clean up
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts

# Install CI Tools
RUN apt-get update \
    && apt-get install -y shellcheck \
    && su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g markdownlint-cli" 2>&1 \
    && pip3 install pre-commit \
    # Clean up
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts

COPY ./scripts /tmp/scripts
RUN apt-get update \
    # Install Go Dev Tools
    && bash /tmp/scripts/install_go_dev_tools.sh \
    # Install Docker CLI
    && bash /tmp/scripts/docker-debian.sh \
    && rm -rf /tmp/scripts

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Uncomment the next line to use go get to install anything else you need
# RUN go get -x <your-dependency-or-tool>

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1

ENTRYPOINT ["/usr/local/share/docker-init.sh"]
CMD ["sleep", "infinity"]
