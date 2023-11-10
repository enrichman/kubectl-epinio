#!/bin/bash

file=$1

# save file to check its pre-edit contents
cat $file > /tmp/before_edit.yaml

# "edit" the file
cat /tmp/edited_user.yaml > $file
