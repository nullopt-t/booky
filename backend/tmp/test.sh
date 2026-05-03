#!/bin/bash

BASE_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "  ORDER ENDPOINT TEST SUITE"
echo "=========================================="
echo ""

# ============================================
# STEP 1: Create test products first
# ============================================
echo -e "${YELLOW}>>> Creating test products...${NC}"

# Product with low stock (for stock error test)
PRODUCT_LOW=$(curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"title":"Low Stock Book","price":1000,"stock":2}' | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

# Product with high stock (for success test)
PRODUCT_HIGH=$(curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"title":"High Stock Book","price":1500,"stock":100}' | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "Low stock product ID: $PRODUCT_LOW"
echo "High stock product ID: $PRODUCT_HIGH"
echo ""

# ============================================
# TEST 1: Invalid JSON Body
# ============================================
echo -e "${YELLOW}TEST 1: Invalid JSON Body (400 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d 'not valid json' | jq .
echo ""
echo ""

# ============================================
# TEST 2: Missing items field
# ============================================
echo -e "${YELLOW}TEST 2: Missing items field (400 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{}' | jq .
echo ""
echo ""

# ============================================
# TEST 3: Empty items array
# ============================================
echo -e "${YELLOW}TEST 3: Empty items array (400 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[]}' | jq .
echo ""
echo ""

# ============================================
# TEST 4: Invalid Product ID (not UUID)
# ============================================
echo -e "${YELLOW}TEST 4: Invalid Product ID - not UUID (400 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"not-a-uuid","quantity":1}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 5: Missing product_id
# ============================================
echo -e "${YELLOW}TEST 5: Missing product_id (400 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"quantity":1}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 6: Missing quantity
# ============================================
echo -e "${YELLOW}TEST 6: Missing quantity (400 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"'$PRODUCT_HIGH'"}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 7: Quantity < 1
# ============================================
echo -e "${YELLOW}TEST 7: Quantity = 0 (400 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"'$PRODUCT_HIGH'","quantity":0}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 8: Quantity > 100
# ============================================
echo -e "${YELLOW}TEST 8: Quantity = 101 (400 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"'$PRODUCT_HIGH'","quantity":101}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 9: Non-existent Product (valid UUID)
# ============================================
echo -e "${YELLOW}TEST 9: Non-existent Product - valid UUID (404 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11","quantity":1}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 10: Insufficient Stock (want 10, have 2)
# ============================================
echo -e "${YELLOW}TEST 10: Insufficient Stock (500 Expected - wrapped as generic)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"'$PRODUCT_LOW'","quantity":10}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 11: Stock edge case - exact amount
# ============================================
echo -e "${YELLOW}TEST 11: Stock edge case - exact amount (201 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"'$PRODUCT_LOW'","quantity":2}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 12: Multiple items - one invalid
# ============================================
echo -e "${YELLOW}TEST 12: Multiple items - one invalid (404/400 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"'$PRODUCT_HIGH'","quantity":1},{"product_id":"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11","quantity":1}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 13: Success - Single item
# ============================================
echo -e "${YELLOW}TEST 13: SUCCESS - Single item (201 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"'$PRODUCT_HIGH'","quantity":5}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 14: Success - Multiple items
# ============================================
echo -e "${YELLOW}TEST 14: SUCCESS - Multiple items (201 Expected)${NC}"
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"'$PRODUCT_HIGH'","quantity":2},{"product_id":"'$PRODUCT_HIGH'","quantity":3}]}' | jq .
echo ""
echo ""

# ============================================
# TEST 15: Race condition - stock race
# ============================================
echo -e "${YELLOW}TEST 15: Race condition test - rapid concurrent requests${NC}"
echo "Sending 5 concurrent orders for 3 stock item..."

# Create product with exactly 3 stock
PRODUCT_RACE=$(curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"title":"Race Test Book","price":500,"stock":3}' | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

# Send 5 concurrent requests, only 3 should succeed
for i in {1..5}; do
  curl -s -X POST "$BASE_URL/orders" \
    -H "Content-Type: application/json" \
    -d '{"items":[{"product_id":"'$PRODUCT_RACE'","quantity":1}]}' &
done
wait
echo "Check which succeeded (201) vs failed (500)"
echo ""

# ============================================
# STEP 3: Verify stock was decremented
# ============================================
echo -e "${YELLOW}>>> Verifying stock decrementation...${NC}"
echo "Product $PRODUCT_HIGH should have 90 stock (started 100, ordered 10):"
curl -s "$BASE_URL/products/$PRODUCT_HIGH" | jq .
echo ""

echo "=========================================="
echo "  TEST SUITE COMPLETE"
echo "=========================================="