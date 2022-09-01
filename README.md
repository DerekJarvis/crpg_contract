This is the repository for my [New Side Quest blog post](https://www.jarvispowered.com/a-new-side-quest/)

## Summary
* Demo dungeon battler game
* Blockchain Technology: Hyperledger Fabric 2.4
* Smart Contract Language: Go 1.17
* Tested on x86-based Ubuntu
* Use [0_init.sh](0_init.sh) to initialize the project, which will:
	* Ensure the Fabric binaries are present
	* Ensure the Firefly binaries are present
	* Ensure the patched Fabconnect docker image is present
	* Restart or Create a development network using Firefly
	* Update the development network to use the patched Fabconnect docker image
* Use [1_reset.sh](1_reset.sh) to start or reset the project, which will:
	* Clean up any previously compiled contract
	* Reset & Restart the dev Firefly stack
	* Package the chaincode
	* Deploy the chaincode
* Use [4_test.sh](4_test.sh) to run tests against the development network which will:
	* Define the contract interface with Firefly
	* Expose the contract through WebAPIs with Firefly
	* Use the WebAPIs to:
		* Create players
		* Create a character
		* Buy Packs
		* Open Packs
		* Create a dungeon
		* List the dungeon
		* Start a dungeon match
		* Score a dungeon match

## Relevant URLs:
### [Firefly Explorer - Home](http://localhost:5000/ui/namespaces/default/home?time=24hours)
This is a dashboard for the local Firefly node which is running the Fabric network. It also runs our Swagger API and an IPFS node.

### [Firefly Explorer - Operations](http://localhost:5000/ui/namespaces/default/activity/operations?time=24hours)
This page will show operations happening on-chain. You will see successes and failures here as result of the tests running.

### [Firefly Sandbox](http://localhost:5108/)
This is a sandbox page that lets you see live events as they happen, as well as define and deploy new contract APIs.

### [Swagger UI](http://127.0.0.1:5000/api/v1/namespaces/default/apis/crpg/api#/)
This link will only work after the contract has been deployed and registered. The start script should take care of this.

This is an automatically generated WebAPI that exposes all of the contract interactions. The test suite uses this WebAPI to do a complete end-to-end test instead of relying on mocks.