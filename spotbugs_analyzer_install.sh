#!/bin/sh
helm install ${JX_SPOTBUGS_ANALYZER_INSTALL_CHART_REPOSITORY}/ext-spotbugs --version ${EXT_VERSION} --set teamNamespace=${EXT_TEAM_NAMESPACE} --name=ext-jacoco
