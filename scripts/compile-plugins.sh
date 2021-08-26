for PLUG in $(find plugins/* -type d -not -name core); do
  NAME=$(basename $PLUG)
  go build -buildmode=plugin -o $PLUG/$NAME.so $PLUG/$NAME.go
done
