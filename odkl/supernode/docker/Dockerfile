FROM @docker.repo@/@docker.image.base@
LABEL from="@docker.repo@/@docker.image.base@" allow-run-rebase=true 
COPY rootfs /
RUN docker/build
