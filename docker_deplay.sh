#!/usr/bin/env bash

REV=$(git log --pretty=format:'%h' -n 1)

NOT_LOGGED_IN='Not logged in. Use "cf login" to log in.'

cf ic info


if [ $? -ne 1 ]
then
    cf ic stop MCP
    cf ic rm MC
        for i in `cf ic images  | grep mcp | awk '{ print $3 } '`; do cf ic rmi $i; done
    cf ic build -t mcp:$REV  .
    cf ic run -P --memory 64 --name MCP registry.ng.bluemix.net/mexicanstraswberry/mcp:$REV

fi


#cf ic run -P -e "CCS_BIND_APP=mcp" --memory 64 --name MCP registry.ng.bluemix.net/mexicanstraswberry/mcep:756d1c4


