FROM fedora

ENV GOPATH=/go/

RUN dnf -y install qpid-proton-c-devel \
	golang \
	glide \
	git \
	make \
	cmake \
	&& mkdir -p /go/{src,bin,pkg}

# make sure the required environment variables have been passed to the build command
ARG EH_TEST_NAMESPACE
ARG EH_TEST_NAME
ARG EH_TEST_SAS_POLICY_NAME
ARG EH_TEST_SAS_POLICY_KEY
ARG EH_TEST_CONSUMER_GROUP

ENV BIN=prometheus-azure-timeseries-adapter

WORKDIR /go/src/github.com/vglafirov/${BIN}
COPY . ./

RUN make bootstrap
RUN make

CMD /go/bin/$BIN