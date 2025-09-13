# .devcontainer/Dockerfile
# Use a minimal Arch Linux base image for rolling releases
FROM archlinux:base-devel

# Avoid interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Update system and install essential development packages
RUN pacman -Syu --noconfirm && \
    pacman -S --noconfirm \
    git \
    go \
    nodejs \
    npm \
    python \
    python-pip \
    ripgrep \
    fzf \
    shellcheck \
    wget \
    curl \
    openssh \
    sudo

# Create a non-root user (recommended for development)
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID
RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

# Install Go tools (example, customize as needed)
RUN go install golang.org/x/tools/gopls@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest

# Install Neovim nightly from official archive
# Note: Check https://github.com/neovim/neovim/releases/nightly for latest URL
RUN curl -LO https://github.com/neovim/neovim/releases/download/nightly/nvim-linux64.tar.gz && \
    tar -C /usr/local --strip-components=1 -xzf nvim-linux64.tar.gz && \
    rm nvim-linux64.tar.gz

# Install AstroVim for the non-root user
# This clones the configuration into the user's home directory
USER $USERNAME
RUN git clone --depth 1 https://github.com/AstroNvim/AstroNvim ~/.config/nvim

# (Optional) Install your custom AstroVim user config
# COPY --chown=$USER_UID:$USER_GID .nvim /home/$USERNAME/.config/nvim/lua/user

# Switch back to root for final steps
USER root

# Clean up to reduce image size
RUN pacman -Scc --noconfirm && rm -rf /var/cache/pacman/pkg/* /tmp/* /var/tmp/*

# Set the default user for the container
USER $USERNAME
