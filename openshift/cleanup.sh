#!/bin/sh

oc delete is/dumbnailer-base is/dumbnailer-s2i is/dumbnailer
oc delete bc/dumbnailer-base bc/dumbnailer-s2i bc/dumbnailer
oc delete dc/dumbnailer svc/dumbnailer route/dumbnailer
