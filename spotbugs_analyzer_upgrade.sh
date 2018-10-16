#!/bin/sh
NAME=ext-spotbugs
helm repo update
helm install ${JX_SPOTBUGS_ANALYZER_UPGRADE_CHART_REPOSITORY}/${NAME} --version ${EXT_VERSION} --set teamNamespace=${EXT_TEAM_NAMESPACE} --name=${NAME}
