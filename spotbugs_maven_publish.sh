#!/bin/sh
jx step collect --provider=${JX_SPOTBUGS_MAVEN_PUBLISH_PROVIDER} --pattern=target/spotbugsXml.xml --classifier=spotbugs
