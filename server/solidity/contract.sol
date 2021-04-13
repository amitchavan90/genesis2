pragma solidity ^0.6.4;

contract smart_contract {
	address private _owner;
	bytes32 name;

	uint256[] batches;

	constructor(bytes32 _name) public {
		_owner = msg.sender;
		name = _name;
	}

	modifier isOwner() {
		require(msg.sender == _owner, "Sender not authorized.");
		_;
	}

	function owner() public view returns (address) {
		return _owner;
	}

	function commitBatch(uint256 batch) public isOwner() {
		batches.push(batch);
	}

	function batchCount() public view isOwner() returns (uint256) {
		return batches.length;
	}
}
