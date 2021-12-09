#!/bin/bash

sed 's/"/\\"/g' contact-center.avsc | tr '\n' ' ' | sed 's/ //g'