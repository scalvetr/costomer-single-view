#!/bin/bash

echo "contact-center-value.avsc"
echo "-------------------------"
sed 's/"/\\"/g' contact-center-value.avsc | tr '\n' ' ' | sed 's/ //g'
echo ""
cat contact-center-value.avsc | tr '\n' ' ' | sed 's/ //g'
echo ""
echo "-------------------------"

echo "contact-center-key.avsc"
echo "-------------------------"
sed 's/"/\\"/g' contact-center-key.avsc | tr '\n' ' ' | sed 's/ //g'
echo ""
cat contact-center-key.avsc | tr '\n' ' ' | sed 's/ //g'
echo ""
echo "-------------------------"
