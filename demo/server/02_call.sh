# Initialize and list tools
(
  cat <<\EOF
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","clientInfo":{"name":"example-client","version":"1.0.0"},"capabilities":{}}}
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}
EOF
) | ./server

(
  cat <<\EOF
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"duckdb","arguments":{"query":"SELECT * FROM 'test_data.csv' WHERE sales > 1500"}}}
EOF
) | ./server
