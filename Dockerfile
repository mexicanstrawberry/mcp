FROM golang:onbuild

# load mcp app VCAP_SERVICES
ENV CCS_BIND_APP=mcp

# expose api port
EXPOSE 5000