# Go Mimic

Mimic is an end to end network simulation tool for the Hive blockchain. It provides a modular, and easy to use localnet with end to end validation and JSON RPC APIs identical to that of Hive mainnet

🚧**Mimic is Under Active Development**🚧

## Capabilities

- Create test accounts using a specified list of public keys and metadata
- Update account metadata or public keys via mock transactions
- Perform real world transactions against a mock environment, respecting signature validation, balance validation, and transaction expiration



### Supported Transactions

Not all transactions are supported on Mimic. However, we have implemented the following transactions deemed most important for application development:

- `custom_json` 🚧
- `account_update` 🚧
- `account_update2` 🚧
- `transfer` 🚧
- `transfer_from_savings` 🚧
- `transfer_to_savings` 🚧


### Supported API Calls


**Hive APIs**
- `account_history_api.get_ops_in_block` 🚧
- `block_api.get_block` 🚧
- `block_api.get_block_range` 🚧
- `condenser_api.broadcast_transaction` 🚧
- `condenser_api.broadcast_transaction_synchronous` 🚧
- `condenser_api.get_dynamic_global_properties` 🚧
- `condenser_api.get_current_median_history_price` 🚧
- `condenser_api.get_reward_fund` 🚧
- `condenser_api.get_withdraw_routes` 🚧
- `condenser_api.get_open_orders` 🚧
- `condenser_api.get_conversion_requests` 🚧
- `condenser_api.get_collateralized_conversion_requests` 🚧
- `condenser_api.get_accounts` 🚧
- `rc_api.find_rc_accounts` 🚧
- `/health` ✅

**Mimic APIs**
- Admin create account / modify keys 🚧
- Admin transaction 🚧
- Admin reset block database



**Virtual Ops**:
- Claim HBD savings 🚧

### Limitations

- Mimic does not currently support creation posts, comments, likes, proof of brain or other social related activties
- Mimic does simulates witness scheduling in ideal conditions, block production may not completely match real world conditions.
- Mimic does not currently support authority delegation. Transactions created using authority delegation will fail
- Mimic may have not implemented all error types or may not be able to simulate all types of error returned from Hive RPC API. 
- Mimic does not simulate resource credits
- Mimic does not simulate Hive Power



