#!/bin/sh
helm repo update
helm install ${JX_SPOTBUGS_ANALYZER_UPGRADE_CHART_REPOSITORY}/ext-spotbugs --version ${EXT_VERSION} --set teamNamespace=${EXT_TEAM_NAMESPACE} --name=ext-jacoco
