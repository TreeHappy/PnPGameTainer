FROM mcr.microsoft.com/devcontainers/base:ubuntu-22.04

# Remove Python if present (optional - base image may have minimal Python)
RUN apt-get update && apt-get remove -y python3* python* && apt-get autoremove -y

# Install Zellij terminal workspace
RUN curl -L https://zellij.dev/launch | bash -s -- --yes

# Install OpenSSH server and sudo
RUN apt-get update && apt-get install -y openssh-server sudo

# Create a non-root user (standard practice for dev containers)
RUN useradd -m -s /bin/bash vscode && \
    echo "vscode ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers.d/vscode

# Install Go
RUN wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz && \
    rm go1.22.0.linux-amd64.tar.gz

# Set up environment variables for Go
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# Create Go workspace directory
RUN mkdir -p $GOPATH/src $GOPATH/bin && chmod -R 777 $GOPATH

# Configure SSH
RUN mkdir /var/run/sshd && \
    echo "PasswordAuthentication no" >> /etc/ssh/sshd_config && \
    echo "PermitUserEnvironment yes" >> /etc/ssh/sshd_config

# Expose SSH port
EXPOSE 22

# Clean up
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

