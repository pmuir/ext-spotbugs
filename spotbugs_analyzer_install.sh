#!/bin/sh
NAME=ext-spotbugs
helm install --repo ${JX_SPOTBUGS_ANALYZER_INSTALL_CHART_REPOSITORY} ${NAME} --version ${EXT_VERSION} --set teamNamespace=${EXT_TEAM_NAMESPACE} --name=${NAME}
