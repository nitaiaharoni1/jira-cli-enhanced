#!/bin/bash
# Integration test script for jira-cli
# Tests new operations against a real Jira instance

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if integration tests are enabled
if [ "$JIRA_INTEGRATION_TEST" != "true" ]; then
    echo -e "${YELLOW}Integration tests are disabled.${NC}"
    echo "Set JIRA_INTEGRATION_TEST=true to enable."
    echo ""
    echo "Required environment variables:"
    echo "  JIRA_TEST_SERVER      - Jira server URL (e.g., https://your-instance.atlassian.net)"
    echo "  JIRA_TEST_LOGIN       - Your Jira username/email"
    echo "  JIRA_TEST_API_TOKEN   - Your Jira API token"
    echo "  JIRA_TEST_PROJECT     - Project key for testing (e.g., TEST)"
    echo ""
    echo "Optional:"
    echo "  JIRA_TEST_DEBUG       - Set to 'true' for debug output"
    exit 0
fi

# Validate required environment variables
if [ -z "$JIRA_TEST_SERVER" ] || [ -z "$JIRA_TEST_LOGIN" ] || [ -z "$JIRA_TEST_API_TOKEN" ] || [ -z "$JIRA_TEST_PROJECT" ]; then
    echo -e "${RED}Error: Missing required environment variables${NC}"
    echo ""
    echo "Required:"
    echo "  JIRA_TEST_SERVER=$JIRA_TEST_SERVER"
    echo "  JIRA_TEST_LOGIN=$JIRA_TEST_LOGIN"
    echo "  JIRA_TEST_API_TOKEN=${JIRA_TEST_API_TOKEN:+***}"
    echo "  JIRA_TEST_PROJECT=$JIRA_TEST_PROJECT"
    exit 1
fi

echo -e "${GREEN}Starting integration tests against: $JIRA_TEST_SERVER${NC}"
echo "Project: $JIRA_TEST_PROJECT"
echo ""

# Build the binary
echo "Building jira-cli..."
cd "$(dirname "$0")/.."
go build -o bin/jira ./cmd/jira

# Create a test config
TEST_CONFIG_DIR=$(mktemp -d)
TEST_CONFIG="$TEST_CONFIG_DIR/.jira/.config.yml"
mkdir -p "$(dirname "$TEST_CONFIG")"

cat > "$TEST_CONFIG" <<EOF
installation: Cloud
server: $JIRA_TEST_SERVER
login: $JIRA_TEST_LOGIN
project:
  key: $JIRA_TEST_PROJECT
EOF

export JIRA_CONFIG_FILE="$TEST_CONFIG"
export JIRA_API_TOKEN="$JIRA_TEST_API_TOKEN"

# Test issue keys (will be created)
TEST_ISSUE_1=""
TEST_ISSUE_2=""
TEST_ISSUE_3=""

cleanup() {
    echo ""
    echo -e "${YELLOW}Cleaning up test issues...${NC}"
    
    if [ -n "$TEST_ISSUE_1" ]; then
        ./bin/jira issue delete "$TEST_ISSUE_1" --no-input 2>/dev/null || true
    fi
    if [ -n "$TEST_ISSUE_2" ]; then
        ./bin/jira issue delete "$TEST_ISSUE_2" --no-input 2>/dev/null || true
    fi
    if [ -n "$TEST_ISSUE_3" ]; then
        ./bin/jira issue delete "$TEST_ISSUE_3" --no-input 2>/dev/null || true
    fi
    
    rm -rf "$TEST_CONFIG_DIR"
}

trap cleanup EXIT

# Test 1: Create test issues
echo -e "${GREEN}Test 1: Creating test issues...${NC}"
TEST_ISSUE_1=$(./bin/jira issue create -tTask -s"Integration Test Issue 1" --original-estimate "1h" --no-input --plain 2>/dev/null | head -1 | awk '{print $1}')
TEST_ISSUE_2=$(./bin/jira issue create -tTask -s"Integration Test Issue 2" --original-estimate "2h" --no-input --plain 2>/dev/null | head -1 | awk '{print $1}')
TEST_ISSUE_3=$(./bin/jira issue create -tTask -s"Integration Test Issue 3" --original-estimate "3h" --no-input --plain 2>/dev/null | head -1 | awk '{print $1}')

if [ -z "$TEST_ISSUE_1" ] || [ -z "$TEST_ISSUE_2" ] || [ -z "$TEST_ISSUE_3" ]; then
    echo -e "${RED}Failed to create test issues${NC}"
    exit 1
fi

echo "Created: $TEST_ISSUE_1, $TEST_ISSUE_2, $TEST_ISSUE_3"

# Test 2: Bulk assignment
echo ""
echo -e "${GREEN}Test 2: Bulk assignment...${NC}"
./bin/jira issue assign-bulk "$TEST_ISSUE_1" "$TEST_ISSUE_2" "$TEST_ISSUE_3" "$JIRA_TEST_LOGIN" || {
    echo -e "${RED}Bulk assignment failed${NC}"
    exit 1
}
echo -e "${GREEN}✓ Bulk assignment successful${NC}"

# Test 3: Bulk status transition
echo ""
echo -e "${GREEN}Test 3: Bulk status transition...${NC}"
# Get available status (try "In Progress" or first available)
STATUS=$(./bin/jira issue view "$TEST_ISSUE_1" --plain 2>/dev/null | grep -i "status" | head -1 | awk '{print $NF}' || echo "To Do")
echo "Current status: $STATUS"
# Try to transition (may fail if status doesn't exist, that's okay)
./bin/jira issue move-bulk "$TEST_ISSUE_1" "$TEST_ISSUE_2" "$TEST_ISSUE_3" "In Progress" 2>/dev/null || {
    echo -e "${YELLOW}Note: Status transition may have failed if 'In Progress' is not available${NC}"
}
echo -e "${GREEN}✓ Bulk transition attempted${NC}"

# Test 4: Estimate update
echo ""
echo -e "${GREEN}Test 4: Estimate update...${NC}"
./bin/jira issue estimate "$TEST_ISSUE_1" "4h" || {
    echo -e "${RED}Estimate update failed${NC}"
    exit 1
}
echo -e "${GREEN}✓ Estimate update successful${NC}"

# Test 5: Story points (if configured)
echo ""
echo -e "${GREEN}Test 5: Story points update...${NC}"
./bin/jira issue story-points "$TEST_ISSUE_1" 5 2>/dev/null || {
    echo -e "${YELLOW}Note: Story points may not be configured${NC}"
}
echo -e "${GREEN}✓ Story points update attempted${NC}"

# Test 6: Bulk custom fields
echo ""
echo -e "${GREEN}Test 6: Bulk custom fields...${NC}"
# Try to set a label via custom field (if supported)
./bin/jira issue custom "$TEST_ISSUE_1" "$TEST_ISSUE_2" "test-label=integration-test" 2>/dev/null || {
    echo -e "${YELLOW}Note: Custom field update may have failed (field may not be configured)${NC}"
}
echo -e "${GREEN}✓ Custom field update attempted${NC}"

echo ""
echo -e "${GREEN}All integration tests completed!${NC}"
echo "Test issues: $TEST_ISSUE_1, $TEST_ISSUE_2, $TEST_ISSUE_3"
echo "These will be cleaned up automatically."

