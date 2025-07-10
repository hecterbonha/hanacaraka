#!/bin/bash

# Test script for Hanacaraka API endpoints
# Make sure the server is running on port 8080 before running this script

BASE_URL="http://localhost:8080"

echo "=== Testing Hanacaraka API ==="
echo

# Test home endpoint
echo "1. Testing home endpoint..."
curl -s "$BASE_URL/"
echo
echo

# Test get all users
echo "2. Testing GET /users..."
curl -s "$BASE_URL/users" | jq '.' 2>/dev/null || curl -s "$BASE_URL/users"
echo
echo

# Test get user by ID
echo "3. Testing GET /users/1..."
curl -s "$BASE_URL/users/1" | jq '.' 2>/dev/null || curl -s "$BASE_URL/users/1"
echo
echo

# Test create new user
echo "4. Testing POST /users (create new user)..."
curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Johnson","email":"alice@example.com"}' | jq '.' 2>/dev/null || curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Johnson","email":"alice@example.com"}'
echo
echo

# Test get all users again to see the new user
echo "5. Testing GET /users (after creating new user)..."
curl -s "$BASE_URL/users" | jq '.' 2>/dev/null || curl -s "$BASE_URL/users"
echo
echo

# Test update user
echo "6. Testing PUT /users/3 (update user)..."
curl -s -X PUT "$BASE_URL/users/3" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Cooper","email":"alice.cooper@example.com"}' | jq '.' 2>/dev/null || curl -s -X PUT "$BASE_URL/users/3" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Cooper","email":"alice.cooper@example.com"}'
echo
echo

# Test health endpoint
echo "7. Testing GET /api/v1/health..."
curl -s "$BASE_URL/api/v1/health" | jq '.' 2>/dev/null || curl -s "$BASE_URL/api/v1/health"
echo
echo

# Test delete user
echo "8. Testing DELETE /users/2..."
curl -s -X DELETE "$BASE_URL/users/2" -w "HTTP Status: %{http_code}\n"
echo

# Test get all users after deletion
echo "9. Testing GET /users (after deletion)..."
curl -s "$BASE_URL/users" | jq '.' 2>/dev/null || curl -s "$BASE_URL/users"
echo
echo

# Test non-existent user
echo "10. Testing GET /users/999 (non-existent user)..."
curl -s "$BASE_URL/users/999" -w "HTTP Status: %{http_code}\n"
echo

echo "=== API Testing Complete ==="
