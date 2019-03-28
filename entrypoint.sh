#!/usr/bin/env bash

set -e
set -u
set -o pipefail

#if [ -n "${PARAMETER_STORE:-}" ]; then
#   export TITAN_API_MID__USUARIOKYRON="$(aws ssm get-parameter --name /${PARAMETER_STORE}/titan_api_mid/kyron/username --output text --query Parameter.Value)"
#   export TITAN_API_MID__CLAVE="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/titan_api_mid/kyron/password --output text --query Parameter.Value)"
#fi

exec ./main "$@"
