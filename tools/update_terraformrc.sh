#!/bin/bash
# Copyright (c) Thomas Geens


# Check if the correct number of parameters is provided
if [ $# -ne 4 ]; then
    echo "Usage: $0 <registry_name> <author_slug> <binary_name> <binary_location>"
    exit 1
fi

REGISTRY_NAME=$1
AUTHOR_SLUG=$2
BINARY_NAME=$3
BINARY_LOCATION=$4
PROVIDER_PATH="$BINARY_LOCATION/"

# Strip "terraform-provider-" prefix from binary name for registry path
REGISTRY_PROVIDER_NAME=${BINARY_NAME#terraform-provider-}

PROVIDER_STRING="\"$REGISTRY_NAME/$AUTHOR_SLUG/$REGISTRY_PROVIDER_NAME\" = \"$PROVIDER_PATH\""

# Check if ~/.terraformrc exists
if [ ! -f ~/.terraformrc ]; then
    echo "Creating new ~/.terraformrc file..."
    cat > ~/.terraformrc << EOT
provider_installation {
  dev_overrides {
    $PROVIDER_STRING
  }
}
EOT
    echo "Created ~/.terraformrc file with provider installation path."
    exit 0
fi

# Terraformrc exists, check if it contains provider_installation section
if grep -q "provider_installation" ~/.terraformrc; then
    # Check if it contains dev_overrides section
    if grep -q "dev_overrides" ~/.terraformrc; then
        # Remove only our specific provider if it exists to avoid duplicates
        grep -v "$REGISTRY_NAME/$AUTHOR_SLUG/$REGISTRY_PROVIDER_NAME" ~/.terraformrc > ~/.terraformrc.tmp

        # Add our provider to existing dev_overrides section
        awk -v provider="$PROVIDER_STRING" '
        {
            print $0;
            if($0 ~ /dev_overrides {/) {
                print "    " provider;
            }
        }' ~/.terraformrc.tmp > ~/.terraformrc.new

        mv ~/.terraformrc.new ~/.terraformrc
        rm ~/.terraformrc.tmp
        echo "Updated existing dev_overrides section in ~/.terraformrc"
    else
        # Add dev_overrides section to provider_installation
        awk -v provider="$PROVIDER_STRING" '
        {
            print $0;
            if($0 ~ /provider_installation {/) {
                print "  dev_overrides {";
                print "    " provider;
                print "  }";
            }
        }' ~/.terraformrc > ~/.terraformrc.tmp

        mv ~/.terraformrc.tmp ~/.terraformrc
        echo "Added dev_overrides section to existing provider_installation in ~/.terraformrc"
    fi
else
    # Add provider_installation block to existing file
    echo "" >> ~/.terraformrc
    cat >> ~/.terraformrc << EOT
provider_installation {
  dev_overrides {
    $PROVIDER_STRING
  }
}
EOT
    echo "Added provider_installation block to existing ~/.terraformrc"
fi

echo "Successfully updated ~/.terraformrc file with provider installation path."
