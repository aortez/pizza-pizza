# Deploys the webapp files to the specified directory.
#!/bin/bash
set -euo pipefail

if [[ $# -eq 0 ]]
  then
    echo "Please specify output directory."
fi

OUTPUT_DIR=$1
echo "Copying files to output dir: $OUTPUT_DIR ..."

cp index.html $OUTPUT_DIR/go.html
cp main.wasm $OUTPUT_DIR
cp wasm_exec.js $OUTPUT_DIR

echo "Done copying files!"
