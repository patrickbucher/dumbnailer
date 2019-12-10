#!/bin/sh

oc delete is/thumbnailer-base is/thumbnailer-s2i is/thumbnailer
oc delete bc/thumbnailer-base bc/thumbnailer-s2i bc/thumbnailer
oc delete dc/thumbnailer svc/thumbnailer route/thumbnailer
