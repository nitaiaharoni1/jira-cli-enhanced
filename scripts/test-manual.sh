#!/bin/bash
# Manual testing script for new operations
# This script helps you test the new operations interactively

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Jira CLI - Manual Testing Script${NC}"
echo ""

# Check if jira is built
if [ ! -f "./bin/jira" ]; then
    echo "Building jira-cli..."
    go build -o bin/jira ./cmd/jira
fi

# Check configuration
if [ -z "$JIRA_API_TOKEN" ]; then
    echo -e "${YELLOW}Warning: JIRA_API_TOKEN not set${NC}"
    echo "Set it with: export JIRA_API_TOKEN='your-token'"
    echo ""
fi

echo -e "${GREEN}Available test operations:${NC}"
echo ""
echo "1. Test Bulk Assignment"
echo "2. Test Bulk Status Transition"
echo "3. Test Estimate Command"
echo "4. Test Story Points Command"
echo "5. Test Bulk Custom Fields"
echo "6. Run All Tests"
echo ""

read -p "Select operation (1-6): " choice

case $choice in
    1)
        echo ""
        echo -e "${BLUE}Testing Bulk Assignment${NC}"
        echo "Enter issue keys (space-separated):"
        read -r issues
        echo "Enter assignee (or 'x' to unassign):"
        read -r assignee
        ./bin/jira issue assign-bulk $issues "$assignee"
        ;;
    2)
        echo ""
        echo -e "${BLUE}Testing Bulk Status Transition${NC}"
        echo "Enter issue keys (space-separated):"
        read -r issues
        echo "Enter target status:"
        read -r status
        ./bin/jira issue move-bulk $issues "$status"
        ;;
    3)
        echo ""
        echo -e "${BLUE}Testing Estimate Command${NC}"
        echo "Enter issue keys (space-separated):"
        read -r issues
        echo "Enter estimate (e.g., '2d 3h'):"
        read -r estimate
        echo "Update remaining estimate? (y/n):"
        read -r remaining
        if [ "$remaining" = "y" ]; then
            ./bin/jira issue estimate $issues "$estimate" --remaining
        else
            ./bin/jira issue estimate $issues "$estimate"
        fi
        ;;
    4)
        echo ""
        echo -e "${BLUE}Testing Story Points Command${NC}"
        echo "Enter issue keys (space-separated):"
        read -r issues
        echo "Enter story points:"
        read -r points
        ./bin/jira issue story-points $issues "$points"
        ;;
    5)
        echo ""
        echo -e "${BLUE}Testing Bulk Custom Fields${NC}"
        echo "Enter issue keys (space-separated):"
        read -r issues
        echo "Enter custom fields (format: FIELD=VALUE,FIELD2=VALUE2):"
        read -r fields
        ./bin/jira issue custom $issues $fields
        ;;
    6)
        echo ""
        echo -e "${BLUE}Running All Tests${NC}"
        echo "This will create test issues and test all operations..."
        echo ""
        read -p "Continue? (y/n): " confirm
        if [ "$confirm" != "y" ]; then
            exit 0
        fi
        
        # Create test issues
        echo "Creating test issues..."
        ISSUE1=$(./bin/jira issue create -tTask -s"Test Bulk 1" --no-input --plain 2>/dev/null | head -1 | awk '{print $1}')
        ISSUE2=$(./bin/jira issue create -tTask -s"Test Bulk 2" --no-input --plain 2>/dev/null | head -1 | awk '{print $1}')
        ISSUE3=$(./bin/jira issue create -tTask -s"Test Bulk 3" --no-input --plain 2>/dev/null | head -1 | awk '{print $1}')
        
        echo "Created: $ISSUE1, $ISSUE2, $ISSUE3"
        
        # Test operations
        echo ""
        echo "Testing bulk assignment..."
        ./bin/jira issue assign-bulk "$ISSUE1" "$ISSUE2" "$ISSUE3" $(./bin/jira me 2>/dev/null | head -1) || true
        
        echo ""
        echo "Testing estimate..."
        ./bin/jira issue estimate "$ISSUE1" "2h" || true
        
        echo ""
        echo "Testing story points..."
        ./bin/jira issue story-points "$ISSUE1" 5 2>/dev/null || echo "Story points may not be configured"
        
        echo ""
        echo "Cleaning up..."
        ./bin/jira issue delete "$ISSUE1" --no-input 2>/dev/null || true
        ./bin/jira issue delete "$ISSUE2" --no-input 2>/dev/null || true
        ./bin/jira issue delete "$ISSUE3" --no-input 2>/dev/null || true
        
        echo -e "${GREEN}All tests completed!${NC}"
        ;;
    *)
        echo "Invalid choice"
        exit 1
        ;;
esac

