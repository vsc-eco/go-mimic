import requests
import json

"""
POST http://localhost:3000 HTTP/1.1
Accept: application/json
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "method": "condenser_api.broadcast_transaction_synchronous",
  "params": [
    {
      "ref_block_num": 1,
      "ref_block_prefix": 0,
      "expiration": "1970-01-01T00:00:00",
      "operations": [],
      "extensions": [],
      "signatures": []
    }
  ],
  "id": 1
}
"""

payload = r"""{
  "ref_block_num": 12345,
  "ref_block_prefix": 987654321,
  "expiration": "2025-07-08T15:30:00",
  "operations": [
    [
      "transfer",
      {
        "from": "alice",
        "to": "bob",
        "amount": "10.000 HIVE",
        "memo": "Payment for services"
      }
    ],
    [
      "vote",
      {
        "voter": "alice",
        "author": "bob",
        "permlink": "my-awesome-post",
        "weight": 10000
      }
    ],
    [
      "comment",
      {
        "parent_author": "",
        "parent_permlink": "hive",
        "author": "alice",
        "permlink": "my-new-post-123",
        "title": "My New Post",
        "body": "This is the content of my post",
        "json_metadata": "{\"tags\":[\"hive\",\"blockchain\"]}"
      }
    ]
  ],
  "extensions": [],
  "signatures": [
    "1f2a3b4c5d6e7f8g9h0i1j2k3l4m5n6o7p8q9r0s1t2u3v4w5x6y7z8a9b0c1d2e3f4g5h6i7j8k9l0m1n2o3p4q5r6s7t8u9v0w1x2y3z4a5b6c7d8e9f0"
  ]
}"""


request_body = json.dumps(
    {
        "jsonrpc": "2.0",
        "method": "condenser_api.broadcast_transaction",
        "params": [json.loads(payload)],
        "id": 1,
    }
)

s = requests.Session()

for i in range(3):
    response = s.post("http://localhost:3000", data=request_body)
    print(response.content)
