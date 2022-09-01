#!/bin/bash
# Make sure fabric binaries exist
if [ -d "fabric" ]; then
	echo "Fabric binaries present"
else
	mkdir fabric
	pushd fabric
	wget https://github.com/hyperledger/fabric/releases/download/v2.4.6/hyperledger-fabric-linux-amd64-2.4.6.tar.gz
	tar -xvf hyperledger-fabric-linux-amd64-2.4.6.tar.gz
	rm hyperledger-fabric-linux-amd64-2.4.6.tar.gz
	popd
fi

# Make sure firefly binaries exist
if [ -d "firefly" ]; then
	echo "Firefly binaries present"
else
	mkdir firefly
	pushd firefly
	wget https://github.com/hyperledger/firefly-cli/releases/download/v1.0.2/firefly-cli_1.0.2_Linux_x86_64.tar.gz
	tar -xvf firefly-cli_1.0.2_Linux_x86_64.tar.gz
	rm firefly-cli_1.0.2_Linux_x86_64.tar.gz
	popd
fi

# Make sure the patched fabconnect exists
if [ -n "$(docker images -q ffi_patch:latest)" ]; then
	echo "Patched fabconnect image already present"
else
	echo "Building patched fabconnect"

	git clone https://github.com/DerekJarvis/firefly-fabconnect.git -b feature-enhanced_type_validation fabconnect-patch
	pushd fabconnect-patch
	docker build -t ffi_patch .
	popd
	rm -r -f fabconnect-patch
fi

# Ensure the dev stack exists
if ./firefly/ff ls | grep "dev"; then
	echo "Stack exists. Resetting it";
	./firefly/ff reset -f dev
else
	echo "Dev Stack does not exist. Creating it."
	./firefly/ff init dev -b fabric 1
fi

# Ensure the docker override file is in place
cp -f docker-compose.override.yml "$(ff info dev | sed -nr 's/.* found at: (.*)docker-compose.yml/\1docker-compose.override.yml/p')"
