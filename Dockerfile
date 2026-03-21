FROM nixos/nix:latest AS base

RUN echo 'experimental-features = nix-command flakes' >> /etc/nix/nix.conf

FROM base AS builder

COPY . /tmp/build
WORKDIR /tmp/build

RUN nix \
    --option filter-syscalls false \
    build

RUN mkdir /tmp/nix-store-closure
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure

FROM base AS dev 

WORKDIR /app

COPY flake.nix flake.lock ./

RUN git config --global --add safe.directory /app
RUN nix --option filter-syscalls false develop -c true

CMD ["nix", "develop", "-c", "air"]

FROM scratch AS prod
WORKDIR /app

COPY --from=builder /tmp/nix-store-closure /nix/store
COPY --from=builder /tmp/build/result /app
CMD ["/app/bin/clicknest"]
