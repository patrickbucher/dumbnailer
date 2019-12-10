#!/bin/sh

curl -X POST -F pdf=@../demo.pdf http://thumbnailer-whatever.192.168.42.154.nip.io/thumbnail > thumbnail.png 
